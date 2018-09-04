package user

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"fmt"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {

	return si.NewGetUsersOK()
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {

	// Tokenのチェック
	if p.Token == "" {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "No Token",
			})
	}

	r1 := repositories.NewUserTokenRepository()
	ent1, err1 := r1.GetByToken(p.Token)
	if err1 != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if ent1 == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "401",
				Message: "Unauthorized Token",
			})
	}

	// Tokenがあった場合の処理
	r := repositories.NewUserRepository()

	ent, err := r.GetByUserID(p.UserID)
	fmt.Println(err)


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
				Message: "User Not Found",
			})
	}

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {

	return si.NewPutProfileOK()
}
