package user

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

//	探すAPI
// 	GET {hostname}/api/1.0/users
func getUsersInternalServerErrorResponse() middleware.Responder {
	return si.NewGetUsersInternalServerError().WithPayload(
		&si.GetUsersInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getUsersUnauthorizedResponse() middleware.Responder {
	return si.NewGetUsersUnauthorized().WithPayload(
		&si.GetUsersUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func getUsersBadRequestResponses() middleware.Responder {
	return si.NewGetUsersBadRequest().WithPayload(
		&si.GetUsersBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

//	ユーザー詳細API
//	GET {hostname}/api/1.0/users/{userID}
func getProfileByUserIDInternalServerErrorResponse() middleware.Responder {
	return si.NewGetProfileByUserIDInternalServerError().WithPayload(
		&si.GetProfileByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getProfileByUserIDNotFoundResponse() middleware.Responder {
	return si.NewGetProfileByUserIDNotFound().WithPayload(
		&si.GetProfileByUserIDNotFoundBody{
			Code:    "404",
			Message: "User Not Found",
		})
}

func getProfileByUserIDUnauthorizedResponse() middleware.Responder {
	return si.NewGetProfileByUserIDUnauthorized().WithPayload(
		&si.GetProfileByUserIDUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func getProfileByUserIDBadRequestResponses() middleware.Responder {
	return si.NewGetProfileByUserIDBadRequest().WithPayload(
		&si.GetProfileByUserIDBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

func getProfileByUserIDLessThan0BadRequestResponses() middleware.Responder {
	return si.NewGetProfileByUserIDBadRequest().WithPayload(
		&si.GetProfileByUserIDBadRequestBody{
			Code:    "400",
			Message: "ID Should Be Bigger Than 0",
		})
}

//	ユーザー情報更新API
//	PUT {hostname}/api/1.0/users/{userID}
func putProfileInternalServerErrorResponse() middleware.Responder {
	return si.NewPutProfileInternalServerError().WithPayload(
		&si.PutProfileInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func putProfileForbiddenResponse() middleware.Responder {
	return si.NewPutProfileForbidden().WithPayload(
		&si.PutProfileForbiddenBody{
			Code:    "403",
			Message: "Forbidden",
		})
}

func putProfileUnauthorizedResponse() middleware.Responder {
	return si.NewPutProfileUnauthorized().WithPayload(
		&si.PutProfileUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func putProfileBadRequestResponse() middleware.Responder {
	return si.NewPutProfileBadRequest().WithPayload(
		&si.PutProfileBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

func putProfileLessThan0BadRequestResponse() middleware.Responder {
	return si.NewPutProfileBadRequest().WithPayload(
		&si.PutProfileBadRequestBody{
			Code:    "400",
			Message: "ID Should Be Bigger Than 0",
		})
}
