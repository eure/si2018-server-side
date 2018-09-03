package user

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/models"
)

// get users
func GetUsers(p si.GetUsersParams) middleware.Responder {

	// get login user by token
	tokenR := repositories.NewUserTokenRepository()

	tokenEnt, err := tokenR.GetByToken(p.Token)

	// Unauthorized
	if err != nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message:  "this token is invalid.",
			})
	}

	// create token from model
	token := tokenEnt.Build()

	// 省くユーザーのid
	var exceptIds []int64

	// 自分のIDを検索に含めない設定をする
	exceptIds = append(exceptIds, token.UserID)

	// ユーザーレポジトリの初期化
	userR := repositories.NewUserRepository()

	// tokenのidからユーザーを取得する
	myUserEnt, err := userR.GetByUserID(token.UserID)

	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// ユーザーモデルを作る
	myUser := myUserEnt.Build()

	// ユーザーマッチレポジトリを初期化する
	userMatchR := repositories.NewUserMatchRepository()

	// マッチしているユーザーを取得する
	matchUserIds, err := userMatchR.FindAllByUserID(myUser.ID)

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

	// ユーザーライクレポジトリを初期化する
	userLikeR := repositories.NewUserLikeRepository()

	// ライクしているユーザーを取得する
	likeUserIds, err := userLikeR.FindLikeAll(myUser.ID)

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

	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 返すユーザーモデルのポインタのスライスを定義する
	var sUsers []*models.User

	// 定義したモデルにマッピングする
	for _, userEnt := range userEnts {
		userModel := userEnt.Build()
		sUsers = append(sUsers, &userModel)
	}

	// 結果を返す
	return si.NewGetUsersOK().WithPayload(sUsers)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	return si.NewGetProfileByUserIDOK()
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	return si.NewPutProfileOK()
}
