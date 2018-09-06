package token

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func getTokenByUserIDThrowInternalServerError(fun string, err error) *si.GetTokenByUserIDInternalServerError {
	return si.NewGetTokenByUserIDInternalServerError().WithPayload(
		&si.GetTokenByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func getTokenByUserIDThrowNotFound(mes string) *si.GetTokenByUserIDNotFound {
	return si.NewGetTokenByUserIDNotFound().WithPayload(
		&si.GetTokenByUserIDNotFoundBody{
			Code:    "404",
			Message: "User Token Not Found: " + mes,
		})
}

func GetTokenByUserID(p si.GetTokenByUserIDParams) middleware.Responder {
	tokenRepo := repositories.NewUserTokenRepository()

	token, err := tokenRepo.GetByUserID(p.UserID)
	if err != nil {
		return getTokenByUserIDThrowInternalServerError("GetByUserID", err)
	}
	if token == nil {
		return getTokenByUserIDThrowNotFound("GetByUserID failed")
	}

	sEnt := token.Build()
	return si.NewGetTokenByUserIDOK().WithPayload(&sEnt)
}
