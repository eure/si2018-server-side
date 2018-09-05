package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	return si.NewPostMessageOK()
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {

	// レポジトリを初期化する
	tokenR := repositories.NewUserTokenRepository()
	usermessageR := repositories.NewUserMessageRepository()

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Token)

	// 401エラー
	if tokenEnt == nil {
		si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code: "401",
				Message: "Token is invalid",
			})
	}

	// 500エラー
	if err != nil {
		si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	// limitのデフォルト値を100に設定する
	var limit int
	if &p.Limit == nil {
		limit = 100
	} else {
		limit = int(*p.Limit)
	}

	messageEnts, err := usermessageR.GetMessages(tokenEnt.UserID, p.UserID, limit, p.Latest, p.Oldest)

	// 500エラー
	if err != nil {
		si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	var messageEntities entities.UserMessages

	messageEntities = messageEnts

	messages := messageEntities.Build()

	return si.NewGetMessagesOK().WithPayload(messages)
}
