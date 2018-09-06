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
	userImageR := repositories.NewUserImageRepository()

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Token)

	// 401エラー
	if tokenEnt == nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message:  "Your token is invalid.",
			})
	}

	// 500エラー
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// マッチ済みのすべてのお相手のUserIdを返す
	userMatchIds, err := userMatchR.FindAllByUserID(tokenEnt.UserID)

	// 500エラー
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// マッチ済みのユーザーを除くいいねが送られたユーザーのいいねを取得する
	userLikeEnts, err := userLikeR.FindGotLikeWithLimitOffset(tokenEnt.UserID, int(p.Limit), int(p.Offset), userMatchIds)

	// 500エラー
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
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
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	var userEntities entities.Users = userEnts

	// ユーザーIDのスライスを取得する
	var userIDs []int64

	for _, userEntity := range userEntities {
		userIDs = append(userIDs, userEntity.ID)
	}

	// 取得したユーザーに紐づく画像を取得する
	userImageEnts, err := userImageR.GetByUserIDs(userIDs)

	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// レスポンス用のスライスに必要な値をマップする
	var likeUserResponses entities.LikeUserResponses

	for _, userLikeEnt := range userLikeEnts {

		var likeUserResponse entities.LikeUserResponse

		// いいねした時間を入れる
		likeUserResponse.LikedAt = userLikeEnt.CreatedAt

		// ユーザーを入れる
		for _, userEnt := range userEnts {

			// ユーザーにimage_uriを入れる
			for _, userImageEnt := range userImageEnts {
				if userEnt.ID == userImageEnt.UserID {
					userEnt.ImageURI = userImageEnt.Path
				}
			}

			if userLikeEnt.UserID == userEnt.ID {
				likeUserResponse.ApplyUser(userEnt)
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
	userR := repositories.NewUserRepository()
	userLikeR := repositories.NewUserLikeRepository()
	userMatchR := repositories.NewUserMatchRepository()

	// 400エラー
	if p.Params.Token == "" {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message:  "Can't find token.",
			})
	}

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)

	// 401エラー
	if tokenEnt == nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message:  "Your token is invalid.",
			})
	}

	// 500エラー
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// ログイン中のユーザーを取得する
	myUserEnt, err := userR.GetByUserID(tokenEnt.UserID)

	// 500エラー
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// いいね先のユーザーを取得する
	partnerUserEnt, err := userR.GetByUserID(p.UserID)

	// 500エラー
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 400エラー
	if myUserEnt.Gender == partnerUserEnt.Gender {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request. Sorry, we are creating genderless society.",
			})
	}

	pastLike, err := userLikeR.GetLikeBySenderIDReceiverID(tokenEnt.UserID, p.UserID)

	// 500エラー
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// いいねは１度しか送れない
	if pastLike != nil {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request. You have already sended a like to this user.",
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
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 相手からのいいねを調べる
	partnerUserLike, err := userLikeR.GetLikeBySenderIDReceiverID(p.UserID, tokenEnt.UserID)

	// 500エラー
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
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
			return si.NewGetLikesInternalServerError().WithPayload(
				&si.GetLikesInternalServerErrorBody{
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
