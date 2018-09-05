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

	// レポジトリを初期化する
	tokenR := repositories.NewUserTokenRepository()
	userLikeR := repositories.NewUserLikeRepository()
	userMatchR := repositories.NewUserMatchRepository()

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)

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

	// いいねを送信する用の構造体を作成する
	userLike := entities.UserLike{
		UserID: tokenEnt.UserID,
		PartnerID: p.UserID,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}

	// いいねを送信する
	err = userLikeR.Create(userLike)

	// 500エラー
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 相手からのいいねを調べる
	partnerUserLike, err := userLikeR.GetLikeBySenderIDReceiverID(p.UserID, tokenEnt.UserID)

	// 500エラー
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 相手からのいいねがある場合、マッチを作成する
	if partnerUserLike != nil {
		userMatchEnt := entities.UserMatch{
			UserID: p.UserID,
			PartnerID: tokenEnt.UserID,
			CreatedAt: strfmt.DateTime(time.Now()),
			UpdatedAt: strfmt.DateTime(time.Now()),
		}

		err := userMatchR.Create(userMatchEnt)

		if err != nil {
			return si.NewGetUsersInternalServerError().WithPayload(
				&si.GetUsersInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	}

	// 結果を返す
	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code: "200",
			Message: "OK",
		})

}
