package token

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetTokenByUserID(p si.GetTokenByUserIDParams) middleware.Responder {
	id := p.UserID
	if id <= 0 {
		return getTokenByUserIDBadRequestResponse()
	}

	r := repositories.NewUserTokenRepository()

	ent, err := r.GetByUserID(id)

	if err != nil {
		return getTokenByUserIDInternalServerErrorResponse()
	}
	if ent == nil {
		return getTokenByUserIDNotFoundResponse()
	}

	sEnt := ent.Build()
	return si.NewGetTokenByUserIDOK().WithPayload(&sEnt)
}
