package token

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetTokenByUserID(p si.GetTokenByUserIDParams) middleware.Responder {
	u := repositories.NewUserTokenRepository()

	// paramsの変数を定義
	paramsUserID := p.UserID

	userToken, err := u.GetByUserID(paramsUserID)

	if err != nil {
		return si.NewGetTokenByUserIDInternalServerError().WithPayload(
			&si.GetTokenByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userToken == nil {
		return si.NewGetTokenByUserIDNotFound().WithPayload(
			&si.GetTokenByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Token Not Found",
			})
	}

	sEnt := userToken.Build()
	return si.NewGetTokenByUserIDOK().WithPayload(&sEnt)
}
