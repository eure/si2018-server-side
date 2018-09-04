package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/models"
	"fmt"
	"encoding/json"
	"strings"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	// TODO: 400エラー
	// TODO: 401エラー

	// TokenからUserIdを取得する
	tokenR         := repositories.NewUserTokenRepository()
	tokenEnt, err  := tokenR.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	//fmt.Println(tokenEnt.UserID)

	// 自分, いいねした人,
	var omitIds []int64
	omitIds = append(omitIds, tokenEnt.UserID)

	// マッチしているユーザを取得する
	machR              := repositories.NewUserMatchRepository()
	matchUserIds, err  := machR.FindAllByUserID(tokenEnt.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	for _, matchUserId := range matchUserIds {
		omitIds = append(omitIds, matchUserId)
	}

	// いいねしているユーザを取得する
	likeR             := repositories.NewUserLikeRepository()
	likeUserIds, err  := likeR.FindLikeAll(tokenEnt.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	for _, likeUserId := range likeUserIds {
		omitIds = append(omitIds, likeUserId)
	}

	// 自分の性別情報の取得
	// omitIds以外のユーザ情報を取得する
	userR       := repositories.NewUserRepository()
	myEnt, err  := userR.GetByUserID(tokenEnt.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	findUsers, err := userR.FindWithCondition(int(p.Limit), int(p.Offset), myEnt.GetOppositeGender(), omitIds)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	// これはentitiesの配列
	//fmt.Println(findUsers)
	//responseData := entities.Users(findUsers).Build()
	var responseData []*models.User
	for _, userEnt := range findUsers {
		userModel    := userEnt.Build()
		responseData  = append(responseData, &userModel)
	}
	fmt.Println(responseData)

	return si.NewGetUsersOK().WithPayload(responseData)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	// Tokenの形式がおかしい -> 401
	if !(strings.HasPrefix(p.Token, "USERTOKEN"))  {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	// Tokenのユーザが存在しない -> 400 Bad Request
	tokenR := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	if tokenEnt == nil{
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// 探しているユーザの情報取得
	userR    := repositories.NewUserRepository()
	findUserEnt, err := userR.GetByUserID(p.UserID)

	// internal server error
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	// not found
	if findUserEnt == nil {
		return si.NewGetTokenByUserIDNotFound().WithPayload(
			&si.GetTokenByUserIDNotFoundBody{
				Code: "404",
				Message: "User Not Found",
			})
	}

	response := findUserEnt.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&response)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	// 現在の自分の情報取得
	userR := repositories.NewUserRepository()
	userEnt, err := userR.GetByUserID(p.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	// 更新処理
	params, _ := p.Params.MarshalBinary()
	json.Unmarshal(params, &userEnt)
	err = userR.Update(userEnt)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	// 更新後のデータを取得
	responseEnt, err := userR.GetByUserID(p.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	responseData := responseEnt.Build()

	return si.NewPutProfileOK().WithPayload(&responseData)
}
