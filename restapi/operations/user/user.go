package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/entities"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	// OFFSET LIMIT Check -> 400
	if p.Offset < 0 || p.Limit <= 0 {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code   : "400",
				Message: "Bad Request",
			})
	}
	tokenR        := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// Tokenのユーザが存在しない -> 401
	if tokenEnt == nil{
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code   : "401",
				Message: "Token Is Invalid",
			})
	}
	// 除外するユーザのIDを作成する
	var omitIds []int64
	omitIds = append(omitIds, tokenEnt.UserID)

	// いいねしている/されているユーザIDを取得する -> これでOK
	// いいねしてくれている人の一覧は `/likes` で取得できるから必要ないと判断
	likeR         := repositories.NewUserLikeRepository()
	likeIds, err  := likeR.FindLikeAll(tokenEnt.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	omitIds = append(omitIds, likeIds...)
	//for _, likeId := range likeIds {
	//	omitIds = append(omitIds, likeId)
	//}

	// 自分の性別情報の取得
	userR       := repositories.NewUserRepository()
	myEnt, err  := userR.GetByUserID(tokenEnt.UserID)
	myGender    := myEnt.GetOppositeGender()
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	// omitIds以外のユーザ情報を取得する
	findUsers, err := userR.FindWithCondition(int(p.Limit), int(p.Offset), myGender, omitIds)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	var tmp entities.Users
	for _, userEnt := range findUsers{
		tmp = append(tmp, userEnt)
	}
	responseData := tmp.Build()
	//var responseData []*models.User
	//for _, userEnt := range findUsers {
	//	userModel    := userEnt.Build()
	//	responseData  = append(responseData, &userModel)
	//}
	return si.NewGetUsersOK().WithPayload(responseData)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	// 同性のプロフィールも確認可能にしている
	tokenR        := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Token)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// Tokenのユーザが存在しない -> 401
	if tokenEnt == nil{
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code   : "401",
				Message: "Token Is Invalid",
			})
	}
	// 探しているユーザの情報取得
	userR            := repositories.NewUserRepository()
	findUserEnt, err := userR.GetByUserID(p.UserID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	myUserEnt, err := userR.GetByUserID(tokenEnt.UserID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// 同性の場合
	if findUserEnt.GetOppositeGender() == myUserEnt.GetOppositeGender() {
		return si.NewGetProfileByUserIDBadRequest().WithPayload(
			&si.GetProfileByUserIDBadRequestBody{
				Code   : "400",
				Message: "Bad Request",
			})
	}
	// ユーザが存在しない -> 404
	if findUserEnt == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code   : "404",
				Message: "User Not Found",
			})
	}

	responseData := findUserEnt.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&responseData)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	tokenR        := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// Tokenのユーザが存在しない -> 401
	if tokenEnt == nil{
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code   : "401",
				Message: "Token Is Invalid",
			})
	}
	// 編集予定のIDと自分のIDが違う -> 403
	if p.UserID != tokenEnt.UserID {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code   : "403",
				Message: "Forbidden",
			})
	}
	// 編集するユーザ情報を作成
	newUserEnt := entities.User{
		ID:             p.UserID,
		Nickname:       p.Params.Nickname,
		ImageURI:       p.Params.ImageURI,
		Tweet:          p.Params.Tweet,
		Introduction:   p.Params.Introduction,
		ResidenceState: p.Params.ResidenceState,
		HomeState:      p.Params.HomeState,
		Education:      p.Params.Education,
		Job:            p.Params.Job,
		AnnualIncome:   p.Params.AnnualIncome,
		Height:         p.Params.Height,
		BodyBuild:      p.Params.BodyBuild,
		MaritalStatus:  p.Params.MaritalStatus,
		Child:          p.Params.Child,
		WhenMarry:      p.Params.WhenMarry,
		WantChild:      p.Params.WantChild,
		Smoking:        p.Params.Smoking,
		Drinking:       p.Params.Drinking,
		Holiday:        p.Params.Holiday,
		HowToMeet:      p.Params.HowToMeet,
		CostOfDate:     p.Params.CostOfDate,
		NthChild:       p.Params.NthChild,
		Housework:      p.Params.Housework,
	}

	userR := repositories.NewUserRepository()
	err    = userR.Update(&newUserEnt)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// 更新後のデータを取得
	responseEnt, err := userR.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	responseData := responseEnt.Build()
	return si.NewPutProfileOK().WithPayload(&responseData)
}
