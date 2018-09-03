package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"strings"
	"strconv"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	return si.NewGetUsersOK()
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	if !(strings.HasPrefix(p.Token, "USERTOKEN")) || !(strings.HasSuffix(p.Token, strconv.FormatInt(p.UserID, 10))) {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	r := repositories.NewUserRepository()
	ent, err := r.GetByUserID(p.UserID)

	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	if ent == nil {
		return si.NewGetTokenByUserIDNotFound().WithPayload(
			&si.GetTokenByUserIDNotFoundBody{
				Code: "404",
				Message: "User Not Found",
			})
	}

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	return si.NewPutProfileOK()
}
