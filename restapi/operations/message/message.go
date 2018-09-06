package message

import (
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	message := repositories.NewUserMessageRepository()
	match := repositories.NewUserMatchRepository()

	// ユーザーID取得用
	token, err := t.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if token == nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// 自分が既にマッチングしている全てのUserIDを取得
	userIDs, err := match.FindAllByUserID(token.UserID)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userIDs == nil {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// マッチしているユーザーかをきちんと確認する
	if CheckMatchUserID(userIDs, p.UserID) {
		// メッセージの値の定義
		m := entities.UserMessage{
			UserID:    token.UserID,
			PartnerID: p.UserID,
			Message:   p.Params.Message,
			CreatedAt: strfmt.DateTime(time.Now()),
			UpdatedAt: strfmt.DateTime(time.Now()),
		}
		// メッセージをインサートする
		err := message.Create(m)
		if err != nil {
			return si.NewPostLikeInternalServerError().WithPayload(
				&si.PostLikeInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	} else {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	message := repositories.NewUserMessageRepository()
	match := repositories.NewUserMatchRepository()

	// ユーザーID取得用
	token, err := t.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if token == nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// matchingしているユーザーの取得
	userIDs, err := match.FindAllByUserID(token.UserID)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userIDs == nil {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// 明示的に宣言
	var m entities.UserMessages
	// マッチしているユーザーかをきちんと確認する
	if CheckMatchUserID(userIDs, p.UserID) {
		// メッセージを取得する
		m, err = message.GetMessages(token.UserID, p.UserID, int(*p.Limit), p.Latest, p.Oldest)
		if err != nil {
			return si.NewGetMessagesInternalServerError().WithPayload(
				&si.GetMessagesInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	} else {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	sEnt := m.Build()
	return si.NewGetMessagesOK().WithPayload(sEnt)
}

//マッチングをしているか確認に使う関数
func CheckMatchUserID(matchID []int64, userID int64) bool {
	for _, m := range matchID {
		if m == userID {
			return true
		}
	}
	return false
}
