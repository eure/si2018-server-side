package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/strfmt"
	"time"
	"strings"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	// Tokenの形式がおかしい -> 401
	if !(strings.HasPrefix(p.Token, "USERTOKEN"))  {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	// Tokenのユーザが存在しない -> 400 Bad Request
	tokenR := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Token)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	if tokenEnt == nil{
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// マッチしている人のidを取得
	matchR := repositories.NewUserMatchRepository()
	matchIds, err := matchR.FindAllByUserID(tokenEnt.UserID)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	// 結局はlikeのユーザ
	likeR := repositories.NewUserLikeRepository()
	// マッチ済みID(matchIds)を除くユーザ情報を取得
	likeEnt, err := likeR.FindGotLikeWithLimitOffset(tokenEnt.UserID, int(p.Limit), int(p.Offset), matchIds)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	// likeのidからユーザ情報取得
	var likedIds []int64
	for _, u := range likeEnt {
		likedIds = append(likedIds, u.UserID)
	}
	// likeしてくれているユーザIDsから各ユーザ情報を取得
	userR := repositories.NewUserRepository()
	responseModels, err := userR.FindByIDs(likedIds)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	var array entities.LikeUserResponses
	for _, u := range responseModels {
		var tmp = entities.LikeUserResponse{}
		tmp.ApplyUser(u)
		array = append(array, tmp)
	}
	responseData := array.Build()

	return si.NewGetLikesOK().WithPayload(responseData)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	// TODO: 既にいいねしていたら？
	// 自分のユーザIDを取得する
	tokenR := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	likeR := repositories.NewUserLikeRepository()
	tmp := entities.UserLike{
		UserID:    tokenEnt.UserID,
		PartnerID: p.UserID,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}
	err = likeR.Create(tmp)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code: "200",
			Message: "OK",
		})
}
