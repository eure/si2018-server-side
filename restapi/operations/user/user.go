package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	return si.NewGetUsersOK()
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	r := repositories.NewUserRepository()

	//token vaildate
	// v := vaildate.ValidateToken(p.UserID, p.Token)
	//
	// if v != true {
	// 	return si.NewGetProfileByUserIDTokenIsInvalid().WithPayload(
	// 		&si.GetProfileByUserIDTokenIsInvalid{
	// 			Code:  "401",
	// 			Message: "Token Is Invalid",
	// 		})
	// }

	ent, err := r.GetByUserID(p.UserID)

	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if ent == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Token Not Found",
			})
	}

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	return si.NewPutProfileOK()
}
