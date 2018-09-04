package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/strfmt"
	"time"
	"strings"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	// Tokenの形式がおかしい -> 401
	if !(strings.HasPrefix(p.Params.Token, "USERTOKEN"))  {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	// Tokenのユーザが存在しない -> 400 Bad Request
	tokenR := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	if tokenEnt == nil{
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// TODO: 送信先がPartnerではない時 -> 403エラー？ -> けど無いから 400エラー？

	messageR := repositories.NewUserMessageRepository()
	// 新しいメッセージの作成
	tmp := entities.UserMessage{
		UserID: tokenEnt.UserID,
		PartnerID: p.UserID,
		Message: p.Params.Message,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}
	err = messageR.Create(tmp)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code: "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	// Tokenの形式がおかしい -> 401
	if !(strings.HasPrefix(p.Token, "USERTOKEN"))  {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	// Tokenのユーザが存在しない -> 400 Bad Request
	tokenR := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	if tokenEnt == nil{
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// TODO: Partnerかどうかのチェックは必要？
	// TODO: それとも、メッセージの取得なのだから既にマッチしていること前提？
	messageR := repositories.NewUserMessageRepository()
	var responseEnts entities.UserMessages
	messages, _ := messageR.GetMessages(tokenEnt.UserID, p.UserID, int(*p.Limit), p.Latest, p.Oldest)
	for _, message := range messages {
		responseEnts = append(responseEnts, message)
	}

	responseData := responseEnts.Build()
	return si.NewGetMessagesOK().WithPayload(responseData)
}
