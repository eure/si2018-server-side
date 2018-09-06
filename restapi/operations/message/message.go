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

	// paramsの変数を定義
	paramsToken := p.Params.Token
	paramsMessage := p.Params.Message
	paramsUserID := p.UserID

	// ユーザーID取得用
	token, err := t.GetByToken(paramsToken)
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

	// メッセージのチェック(空白ではないか，1万字いないか)
	if paramsMessage == "" {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	} else if len(paramsMessage) > 10000 {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// マッチしているユーザーかをきちんと確認する
	if CheckMatchUserID(userIDs, paramsUserID) {
		// メッセージの値の定義
		m := entities.UserMessage{
			UserID:    token.UserID,
			PartnerID: paramsUserID,
			Message:   paramsMessage,
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

	// paramsの変数を定義
	paramsToken := p.Token
	paramsUserID := p.UserID
	paramsLimit := p.Limit
	paramsLatest := p.Latest
	paramsOldest := p.Oldest

	// ユーザーID取得用
	token, err := t.GetByToken(paramsToken)
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

	// limitが20になっているかをvalidation
	if *paramsLimit != int64(20) {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// マッチしているユーザーかをきちんと確認する
	if CheckMatchUserID(userIDs, paramsUserID) {
		// メッセージを取得する
		m, err = message.GetMessages(token.UserID, paramsUserID, int(*paramsLimit), paramsLatest, paramsOldest)
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
