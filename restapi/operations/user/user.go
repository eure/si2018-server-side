package user

import (
	"encoding/json"
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

	userEnt, err := r.GetByUserID(p.UserID)

	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userEnt == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Profile Not Found",
			})
	}

	user := userEnt.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&user)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	r := repositories.NewUserRepository()

	userEnt, _:= r.GetByUserID(p.UserID)

	// paramsをjsonに変換
	params, _ := p.Params.MarshalBinary()
	// userEntにjsonに変換したparamを入れる
	json.Unmarshal(params, &userEnt)

	err := r.Update(userEnt)

	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	return si.NewPutProfileOK()
}
