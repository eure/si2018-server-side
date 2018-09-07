package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMessagesInternalServerErrorResponse() middleware.Responder {
	return si.NewGetMessagesInternalServerError().WithPayload(
		&si.GetMessagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func GetMessageUnauthorizedResponse() middleware.Responder {
	return si.NewGetMessagesUnauthorized().WithPayload(
		&si.GetMessagesUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func PostMessageOKResponse(message string) middleware.Responder {
	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: message,
		})
}

func PostMessageBadREquestResponse(message string) middleware.Responder {
	return si.NewPostMessageBadRequest().WithPayload(
		&si.PostMessageBadRequestBody{
			Code:    "400",
			Message: message,
		})
}

func PostMessageUnauthorizedResponse() middleware.Responder {
	return si.NewPostMessageUnauthorized().WithPayload(
		&si.PostMessageUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func PostMessageInternalServerErrorResponse() middleware.Responder {
	return si.NewPostMessageInternalServerError().WithPayload(
		&si.PostMessageInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}
