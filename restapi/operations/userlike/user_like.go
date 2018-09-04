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
	// Tokenのユーザが存在しない -> 401
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
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code: "401",
				Message: "Token Is Invalid",
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
	// Tokenの形式がおかしい -> 401
	if !(strings.HasPrefix(p.Params.Token, "USERTOKEN"))  {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	// Tokenのユーザが存在しない -> 401
	tokenR := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	if tokenEnt == nil{
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code: "401",
				Message: "Token Is Invalid",
			})
	}

	// TODO: 相手が異性かどうか
	userR := repositories.NewUserRepository()
	toUserEnt, err := userR.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	myUserEnt, err := userR.GetByUserID(tokenEnt.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	// 同性の場合
	if toUserEnt.GetOppositeGender() == myUserEnt.GetOppositeGender(){
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// TODO: 既にいいねしていたら？
	likeR := repositories.NewUserLikeRepository()
	check, err := likeR.GetLikeBySenderIDReceiverID(myUserEnt.ID, toUserEnt.ID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	// 既にいいねしていた
	if check != nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// TODO: いいねを作成
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
	// Created Like

	// TODO: いいね同士でマッチング
	// 相手から既にいいねをもらっている
	check, err = likeR.GetLikeBySenderIDReceiverID(toUserEnt.ID, myUserEnt.ID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	if check != nil {
		matchR := repositories.NewUserMatchRepository()
		tmp := entities.UserMatch{
			UserID:myUserEnt.ID,
			PartnerID:toUserEnt.ID,
			CreatedAt:strfmt.DateTime(time.Now()),
			UpdatedAt:strfmt.DateTime(time.Now()),
		}
		err := matchR.Create(tmp)
		if err != nil {
			return si.NewPostLikeInternalServerError().WithPayload(
				&si.PostLikeInternalServerErrorBody{
					Code: "500",
					Message: "Internal Server Error",
				})
		}
		// Created Match
	}

	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code: "200",
			Message: "OK",
		})
}
