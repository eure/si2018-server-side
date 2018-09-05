package message

import (
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	rm := repositories.NewUserMessageRepository()
	rmt := repositories.NewUserMatchRepository()
	rt := repositories.NewUserTokenRepository()
	// r := repositories.NewUserRepository()
	t, err := rt.GetByToken(p.Params.Token)
	// トークン認証
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByToken failed: " + err.Error(),
			})
	}
	if t == nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗): GetByToken failed",
			})
	}
	// マッチしているかの確認
	match, err := rmt.Get(p.UserID, t.UserID)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: Get failed: " + err.Error(),
			})
	}
	if match == nil {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request: Get failed",
			})
	}
	now := strfmt.DateTime(time.Now())
	mes := entities.UserMessage{
		UserID:    t.UserID,
		PartnerID: p.UserID,
		Message:   p.Params.Message,
		UpdatedAt: now,
		CreatedAt: now}
	rm.Create(mes)
	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	rm := repositories.NewUserMessageRepository()
	rt := repositories.NewUserTokenRepository()
	r := repositories.NewUserRepository()
	t, err := rt.GetByToken(p.Token)
	// トークン認証
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByToken failed: " + err.Error(),
			})
	}
	if t == nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗): GetByToken failed",
			})
	}
	// p.UserID の実在確認
	user, err := r.GetByUserID(p.UserID)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByUserID failed: " + err.Error(),
			})
	}
	if user == nil {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request: GetByUserID failed",
			})
	}
	// p.Limit のデフォルトは 100 (restapi/summerintern/get_messages_parameters.go)
	var limit int
	if p.Limit == nil {
		limit = 100
	} else {
		limit = int(*p.Limit)
	}
	mes, err := rm.GetMessages(p.UserID, t.UserID, limit, p.Latest, p.Oldest)
	ent := entities.UserMessages(mes)
	return si.NewGetMessagesOK().WithPayload(ent.Build())
}
