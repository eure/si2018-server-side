package user

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	userRepo := repositories.NewUserRepository()
	tokenRepo := repositories.NewUserTokenRepository()
	likeRepo := repositories.NewUserLikeRepository()

	tokenOwner, err := tokenRepo.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:"401",
				Message:"Authorization Required",
			})
	}

	user, err := userRepo.GetByUserID(tokenOwner.UserID)
	if err != nil {
		
	}
	
	user.GetOppositeGender()
	
	users, err := likeRepo.FindLikeAll(tokenOwner.UserID)
	if err != nil {
		
	}
	
	ent, err := userRepo.FindWithCondition(int(p.Limit),int(p.Offset),user.GetOppositeGender(),users)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:"500",
				Message:"Internal Server Error",
			})

	}

	//if err == nil {
	//	return si.NewGetUsersBadRequest().WithPayload(
	//		&si.GetUsersBadRequestBody{
	//			Code:"400",
	//			Message:"Bad Request",
	//		})
	//}

	ent2 := entities.Users(ent)
	sEnt := ent2.Build()
	return si.NewGetUsersOK().WithPayload(sEnt)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	r := repositories.NewUserRepository()

	ent, err := r.GetByUserID(p.UserID)
	if err != nil{
	}

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	return si.NewPutProfileOK()
}
