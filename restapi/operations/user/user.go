package user

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
)

// get users
func GetUsers(p si.GetUsersParams) middleware.Responder {

	// ユーザートークンレポジトリを初期化する
	tokenR := repositories.NewUserTokenRepository()
	userR := repositories.NewUserRepository()
	userLikeR := repositories.NewUserLikeRepository()
	userImageR := repositories.NewUserImageRepository()
	userMatchR := repositories.NewUserMatchRepository()

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Token)

	// 401エラー
	if tokenEnt == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message:  "Your token is invalid.",
			})
	}

	// 500エラー
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 省くユーザーのid
	var exceptIds []int64

	// 自分のIDを検索に含めない設定をする
	exceptIds = append(exceptIds, tokenEnt.UserID)

	// tokenのidからユーザーを取得する
	myUserEnt, err := userR.GetByUserID(tokenEnt.UserID)

	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// マッチしているユーザーを取得する
	matchUserIds, err := userMatchR.FindAllByUserID(tokenEnt.UserID)

	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// マッチしているユーザーを除く設定をする
	for _, matchUserId := range matchUserIds {
		exceptIds = append(exceptIds, matchUserId)
	}

	// ライクしているユーザーを取得する
	likeUserIds, err := userLikeR.FindLikeAll(myUserEnt.ID)

	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// ライクしているユーザーを除く設定をする
	for _, likeUserId := range likeUserIds {
		exceptIds = append(exceptIds, likeUserId)
	}

	// 指定の状態からユーザーを複数取得する
	userEnts, err := userR.FindWithCondition(int(p.Limit), int(p.Offset), myUserEnt.GetOppositeGender(), exceptIds)

	// 500エラー
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 取得したユーザーIDの配列を作成する
	var userIDs []int64

	for _, userEnt := range userEnts {
		userIDs = append(userIDs, userEnt.ID)
	}

	// 取得したユーザーに紐づくimage_uriを取得する
	userImageEnts, err := userImageR.GetByUserIDs(userIDs)

	// 500エラー
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// userモデルにマップする
	var userEntities entities.Users = userEnts

	users := userEntities.Build()

	// user の image_uriに値を入れる
	for _, user := range users {
		for _, userImageEnt := range userImageEnts {
			if user.ID == userImageEnt.UserID {
				user.ImageURI = userImageEnt.Path
			}
		}
	}

	// 結果を返す
	return si.NewGetUsersOK().WithPayload(users)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {

	// ユーザートークンレポジトリを初期化する
	tokenR := repositories.NewUserTokenRepository()
	userR := repositories.NewUserRepository()
	userImageR := repositories.NewUserImageRepository()

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Token)

	// 401エラー
	if tokenEnt == nil {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message:  "Your token is invalid.",
			})
	}

	// 500エラー
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 指定したユーザーを取得する
	myUserEnt, err := userR.GetByUserID(p.UserID)

	// 400エラー
	if myUserEnt == nil {
		return si.NewGetProfileByUserIDBadRequest().WithPayload(
			&si.GetProfileByUserIDBadRequestBody{
				Code: "400",
				Message: "Bad Request. Can't get a user.",
			})
	}

	// 500エラー
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// ユーザー紐づくimage_uriを取得する
	myUserImageEnt, err := userImageR.GetByUserID(tokenEnt.UserID)

	// 500エラー
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// ユーザーモデルを作る
	myUser := myUserEnt.Build()

	// imageUriを入れる
	if myUserImageEnt != nil {
		myUser.ImageURI = myUserImageEnt.Path
	}

	// 結果を返す
	return si.NewGetProfileByUserIDOK().WithPayload(&myUser)

}

func PutProfile(p si.PutProfileParams) middleware.Responder {

	// レポジトリを初期化する
	tokenR := repositories.NewUserTokenRepository()
	userR := repositories.NewUserRepository()
	userImageR := repositories.NewUserImageRepository()

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)

	// 400 エラー.
	// ログインしているユーザーと更新したいユーザーが異なるとき
	if tokenEnt.UserID != p.UserID {
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code:    "400",
				Message:  "Bad Request. You can put only your data.",
			})
	}

	// 401エラー
	if tokenEnt == nil {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message:  "Your token is invalid.",
			})
	}

	// 500エラー
	if err != nil {
		if err != nil {
			return si.NewPutProfileInternalServerError().WithPayload(
				&si.PutProfileInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	}

	// ユーザーエンティティにリクエストをマップする
	userEnt := entities.User{}

	userEnt.ID = p.UserID
	userEnt.AnnualIncome = p.Params.AnnualIncome
	userEnt.BodyBuild = p.Params.BodyBuild
	userEnt.Child = p.Params.Child
	userEnt.CostOfDate = p.Params.CostOfDate
	userEnt.Drinking = p.Params.Drinking
	userEnt.Education = p.Params.Education
	userEnt.Height = p.Params.Height
	userEnt.Holiday = p.Params.Holiday
	userEnt.HomeState = p.Params.HomeState
	userEnt.Housework = p.Params.Housework
	userEnt.HowToMeet = p.Params.HowToMeet
	userEnt.ImageURI = p.Params.ImageURI
	userEnt.Introduction = p.Params.Introduction
	userEnt.Job = p.Params.Job
	userEnt.MaritalStatus = p.Params.MaritalStatus
	userEnt.Nickname = p.Params.Nickname
	userEnt.NthChild = p.Params.NthChild
	userEnt.ResidenceState = p.Params.ResidenceState
	userEnt.Smoking = p.Params.Smoking
	userEnt.Tweet = p.Params.Tweet
	userEnt.WantChild = p.Params.WantChild
	userEnt.WhenMarry = p.Params.WhenMarry

	// ユーザーを更新する
	err = userR.Update(&userEnt)

	// 500エラー
	if err != nil {
		if err != nil {
			return si.NewPutProfileInternalServerError().WithPayload(
				&si.PutProfileInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	}

	// 更新後のユーザーを取得する
	myUpdatedUserEnt, err := userR.GetByUserID(p.UserID)

	// 500エラー
	if err != nil {
		if err != nil {
			return si.NewPutProfileInternalServerError().WithPayload(
				&si.PutProfileInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	}

	// ユーザーのimage_uriを取得する
	myUserImageEnt, err := userImageR.GetByUserID(myUpdatedUserEnt.ID)

	// 500エラー
	if err != nil {
		if err != nil {
			return si.NewPutProfileInternalServerError().WithPayload(
				&si.PutProfileInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	}

	// モデルにマッピングする
	myUpdatedUser := myUpdatedUserEnt.Build()
	
	if myUserImageEnt != nil {
		myUpdatedUser.ImageURI = myUserImageEnt.Path
	}

	// 結果を返す
	return si.NewPutProfileOK().WithPayload(&myUpdatedUser)
}
