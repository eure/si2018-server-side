package message

import (
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func postMessageThrowInternalServerError(fun string, err error) *si.PostMessageInternalServerError {
	return si.NewPostMessageInternalServerError().WithPayload(
		&si.PostMessageInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func postMessageThrowUnauthorized(mes string) *si.PostMessageUnauthorized {
	return si.NewPostMessageUnauthorized().WithPayload(
		&si.PostMessageUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func postMessageThrowBadRequest(mes string) *si.PostMessageBadRequest {
	return si.NewPostMessageBadRequest().WithPayload(
		&si.PostMessageBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func PostMessage(p si.PostMessageParams) middleware.Responder {
	var err error
	messageRepo := repositories.NewUserMessageRepository()
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		t, err := tokenRepo.GetByToken(p.Params.Token)
		// トークン認証
		if err != nil {
			return postMessageThrowInternalServerError("GetByToken", err)
		}
		if t == nil {
			return postMessageThrowUnauthorized("GetByToken failed")
		}
		id = t.UserID
	}
	// マッチしているかの確認
	{
		matchRepo := repositories.NewUserMatchRepository()
		match, err := matchRepo.Get(p.UserID, id)
		if err != nil {
			return postMessageThrowInternalServerError("Get", err)
		}
		if match == nil {
			return postMessageThrowBadRequest("Get failed")
		}
	}
	// メッセージを書きこむ
	now := strfmt.DateTime(time.Now())
	mes := entities.UserMessage{
		UserID:    id,
		PartnerID: p.UserID,
		Message:   p.Params.Message,
		UpdatedAt: now,
		CreatedAt: now}
	err = messageRepo.Create(mes)
	if err != nil {
		return postMessageThrowInternalServerError("Create", err)
	}
	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

func getMessagesThrowInternalServerError(fun string, err error) *si.GetMessagesInternalServerError {
	return si.NewGetMessagesInternalServerError().WithPayload(
		&si.GetMessagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func getMessagesThrowUnauthorized(mes string) *si.GetMessagesUnauthorized {
	return si.NewGetMessagesUnauthorized().WithPayload(
		&si.GetMessagesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func getMessagesThrowBadRequest(mes string) *si.GetMessagesBadRequest {
	return si.NewGetMessagesBadRequest().WithPayload(
		&si.GetMessagesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	var err error
	messageRepo := repositories.NewUserMessageRepository()
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Token)
		if err != nil {
			return getMessagesThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return getMessagesThrowUnauthorized("GetByToken failed")
		}
		id = token.UserID
	}
	// p.Limit のデフォルトは 100 (restapi/summerintern/get_messages_parameters.go)
	var limit int
	if p.Limit == nil {
		limit = 100
	} else {
		limit = int(*p.Limit)
	}
	// メッセージの取得
	message, err := messageRepo.GetMessages(p.UserID, id, limit, p.Latest, p.Oldest)
	if err != nil {
		return getMessagesThrowInternalServerError("GetMessages", err)
	}
	ent := entities.UserMessages(message)
	return si.NewGetMessagesOK().WithPayload(ent.Build())
}
