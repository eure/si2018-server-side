package message

import (
	"github.com/go-openapi/runtime/middleware"

	t "github.com/eure/si2018-server-side/restapi/operations/token"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func ValidatePostMessageParams(msg, token string) middleware.Responder {
	if len(msg) == 0 || len(token) == 0 {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	me, err := t.GetUserByToken(token)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if me == nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Unaothorized",
			})
	}

	return nil
}

func ValidateGetMessagesParams(token string) middleware.Responder {
	if len(token) == 0 {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	me, err := t.GetUserByToken(token)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if me == nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Unaothorized",
			})
	}

	return nil
}
