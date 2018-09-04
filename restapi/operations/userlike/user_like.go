package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {

	// レポジトリを初期化する
	tokenR := repositories.NewUserTokenRepository()
	userMatchR := repositories.NewUserMatchRepository()
	userLikeR := repositories.NewUserLikeRepository()
	userR := repositories.NewUserRepository()

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

	// マッチ済みのすべてのお相手のUserIdを返す
	userMatchIds, err := userMatchR.FindAllByUserID(tokenEnt.UserID)

	// 500エラー
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// マッチ済みのユーザーを除くいいねが送られたユーザーのいいねを取得する
	userLikeEnts, err := userLikeR.FindGotLikeWithLimitOffset(tokenEnt.UserID, int(p.Limit), int(p.Offset), userMatchIds)

	// 500エラー
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// いいねを送ったユーザーのUserIDの配列を作成する
	var userLikeUserIds []int64

	for _, userLikeEnt := range userLikeEnts {

		userLikeUserIds = append(userLikeUserIds, userLikeEnt.UserID)

	}

	// いいねを送ったユーザーを取得する
	userEnts, err := userR.FindByIDs(userLikeUserIds)

	// 500エラー
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// レスポンス用のスライスに必要な値をマップする
	var likeUserResponses entities.LikeUserResponses

	for _, userLikeEnt := range userLikeEnts {

		var likeUserResponse entities.LikeUserResponse

		likeUserResponse.LikedAt = userLikeEnt.CreatedAt

		for _, userEnt := range userEnts {

			if userLikeEnt.UserID == userEnt.ID {
				likeUserResponse.User = userEnt
			}
		}

		likeUserResponses = append(likeUserResponses, likeUserResponse)

	}

	// モデルのスライスにする
	likeUser := likeUserResponses.Build()

	// 結果を返す
	return si.NewGetLikesOK().WithPayload(likeUser)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	return si.NewPostLikeOK()
}
