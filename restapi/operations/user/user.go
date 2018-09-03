package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	// r := repositories.NewUserRepository()
	//
	// usersEnt, err := r.FindWithCondition(p.Limit, p.Offset, p.Token)
	//
	// // FindWithCondition(limit, offset int, gender string, ids []int64)
	//
	// for i, userEnt := range usersEnt {
	// 	users := userEnt.Build()
	// }
	//
	// return si.NewGetUsersOK().WithPayload(&users)
	return si.NewGetUsersOK()
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	r := repositories.NewUserRepository()

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
				Message: "User Profile Not Found",
			})
	}

	userEnt := ent.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&userEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	return si.NewPutProfileOK()
}
