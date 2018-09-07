package token

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func getTokenByUserIDInternalServerErrorResponse() middleware.Responder {
	return si.NewGetTokenByUserIDInternalServerError().WithPayload(
		&si.GetTokenByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getTokenByUserIDNotFoundResponse() middleware.Responder {
	return si.NewGetTokenByUserIDNotFound().WithPayload(
		&si.GetTokenByUserIDNotFoundBody{
			Code:    "404",
			Message: "User Token Not Found",
		})
}
