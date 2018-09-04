package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/strfmt"
	"time"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	// TokenからユーザId取得
	tokenR := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	// マッチしている人のidを取得
	matchR := repositories.NewUserMatchRepository()
	matchIds, err := matchR.FindAllByUserID(tokenEnt.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	// 結局はlikeのユーザ
	likeR := repositories.NewUserLikeRepository()
	likeEnt, err := likeR.FindGotLikeWithLimitOffset(tokenEnt.UserID, int(p.Limit), int(p.Offset), matchIds)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	//fmt.Println(likeEnt)
	// likeのidからユーザ情報取得
	var matchUserIds []int64
	for _, u := range likeEnt {
		matchUserIds = append(matchUserIds, u.UserID)
	}
	userR := repositories.NewUserRepository()
	responseModels, err := userR.FindByIDs(matchUserIds)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
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

	return si.NewGetLikesOK().WithPayload(array.Build())
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
