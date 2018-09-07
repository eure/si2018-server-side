package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/strfmt"
	"time"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	if p.Params.Token == "" {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code   : "400",
				Message: "Bad Request",
			})
	}
	// パラメータのmessageが存在しない可能性があった -> 400
	if p.Params.Message == "" {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code   : "400",
				Message: "Bad Request",
			})
	}
	tokenR        := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// Tokenのユーザが存在しない -> 401
	if tokenEnt == nil{
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code   : "401",
				Message: "Token Is Invalid",
			})
	}

	// 送信先がPartnerではない時 -> 400
	matchR         := repositories.NewUserMatchRepository()
	matchData, err := matchR.Get(tokenEnt.UserID, p.UserID)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	if matchData == nil {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code   : "400",
				Message: "Bad Request",
			})
	}

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
			Code   : "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	// LIMITが0以下 -> 400
	if *p.Limit <= 0 {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code   : "400",
				Message: "Bad Request",
			})
	}
	tokenR        := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// Tokenのユーザが存在しない -> 401
	if tokenEnt == nil{
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code   : "401",
				Message: "Token Is Invalid",
			})
	}
	// 時間比較
	if p.Latest != nil && p.Oldest != nil {
		if time.Time(*p.Latest).Before(time.Time(*p.Oldest)) {
			return si.NewGetMessagesBadRequest().WithPayload(
				&si.GetMessagesBadRequestBody{
					Code   : "400",
					Message: "Bad Request",
				})
		}
	}

	// Partnerかどうかのチェック
	matchR         := repositories.NewUserMatchRepository()
	matchData, err := matchR.Get(tokenEnt.UserID, p.UserID)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// Partnerでない -> 400
	if matchData == nil {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code   : "400",
				Message: "Bad Request",
			})
	}

	messageR    := repositories.NewUserMessageRepository()
	messages, _ := messageR.GetMessages(tokenEnt.UserID, p.UserID, int(*p.Limit), p.Latest, p.Oldest)
	var responseEnt entities.UserMessages
	for _, message := range messages {
		responseEnt = append(responseEnt, message)
	}

	responseData := responseEnt.Build()
	return si.NewGetMessagesOK().WithPayload(responseData)
}
