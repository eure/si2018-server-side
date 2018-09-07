package user

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsersInteralServerErrorResponse(message string) middleware.Responder {
	return si.NewGetUsersInternalServerError().WithPayload(
		&si.GetUsersInternalServerErrorBody{
			Code:    "500",
			Message: message,
		})
}

func GetUsersUnauthorizedResponse(message string) middleware.Responder {
	return si.NewGetUsersUnauthorized().WithPayload(
		&si.GetUsersUnauthorizedBody{
			Code:    "401",
			Message: message,
		})
}

func GetUserProfileByUserIDInternalServerErrorResponse(message string) middleware.Responder {
	return si.NewGetProfileByUserIDInternalServerError().WithPayload(
		&si.GetProfileByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: message,
		})
}

func GetUserProfileByUserIDNotFoundResponse(message string) middleware.Responder {
	return si.NewGetProfileByUserIDNotFound().WithPayload(
		&si.GetProfileByUserIDNotFoundBody{
			Code:    "404",
			Message: message,
		})
}

func GetUserProfileByUserIDUnauthorizeResponse(message string) middleware.Responder {
	return si.NewGetProfileByUserIDUnauthorized().WithPayload(
		&si.GetProfileByUserIDUnauthorizedBody{
			Code:    "401",
			Message: message,
		})
}

func PutProfileInternalServerErrorResponse() middleware.Responder {
	return si.NewPutProfileInternalServerError().WithPayload(
		&si.PutProfileInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func PutProfileUnauthorizedResponse() middleware.Responder {
	return si.NewPutProfileUnauthorized().WithPayload(
		&si.PutProfileUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func PutProfileForbiddenResponse() middleware.Responder {
	return si.NewPutProfileForbidden().WithPayload(
		&si.PutProfileForbiddenBody{
			Code:    "403",
			Message: "Forbidden",
		})
}
