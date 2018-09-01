package user

import (
	middleware "github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func ValidateGetUsers(limit, offset int64, t string) middleware.Responder {
	if limit == 0 {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	if len(t) == 0 {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	return nil
}

func ValidateGetProfileByUserID(t string) middleware.Responder {
	if len(t) == 0 {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	return nil
}

func ValidatePutProfile(t string) middleware.Responder {
	if len(t) == 0 {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	return nil
}
