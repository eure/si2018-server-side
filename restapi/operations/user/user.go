package user

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	ur := repositories.NewUserRepository()
	lr := repositories.NewUserLikeRepository()

	ent, _ := ur.GetByToken(p.Token)
	likes, _ := lr.FindLikeAll(ent.ID)

	var ents entities.Users
	ents, err := ur.FindWithCondition(int(p.Limit), int(p.Offset), ent.GetOppositeGender(), likes)

	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	sEnts := ents.Build()
	
	return si.NewGetUsersOK().WithPayload(sEnts)
}


// 詳細
func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	ur := repositories.NewUserRepository()
	tr := repositories.NewUserTokenRepository()
	ent, err := ur.GetByUserID(p.UserID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
			Code: "500",
			Message: "Internal Server Error",
			})
	}
	if ent == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code: "404",
				Message: "UserID Not Found",
			})
	}
	t, _ := tr.GetByUserID(p.UserID)
	if p.Token != t.Token {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code: "401",
				Message: "Unauthorized",
			})
	}
	
	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

// update
func PutProfile(p si.PutProfileParams) middleware.Responder {
	return si.NewPutProfileOK()
}
