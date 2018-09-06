package user

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
)

//getUser
func GetUsers(p si.GetUsersParams) middleware.Responder {
	tr := repositories.NewUserTokenRepository()
	ur := repositories.NewUserRepository()
	lr := repositories.NewUserLikeRepository()

	token, err := tr.GetByToken(p.Token)
	if err != nil {
		return si.NewGetTokenByUserIDInternalServerError().WithPayload(
			&si.GetTokenByUserIDInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error (in token)",
			})
	}

	if token == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code: "401",
				Message: "Unauthorized (token not found)",
			})
	}

	usr, err := ur.GetByUserID(token.UserID)
	if err != nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	if usr == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code: "404",
				Message: "User not found",
			})
	}

	likes, err := lr.FindLikeAll(usr.ID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	if p.Limit <= 0 {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	if p.Offset < 0 {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	var ents entities.Users
	ents, err = ur.FindWithCondition(int(p.Limit), int(p.Offset), usr.GetOppositeGender(), likes)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	if ents == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	sEnts := ents.Build()
	
	return si.NewGetUsersOK().WithPayload(sEnts)
}


// 詳細
func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	ur := repositories.NewUserRepository()
	tr := repositories.NewUserTokenRepository()
	
	token, err := tr.GetByToken(p.Token)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error (in get token)",
			})
	}

	if token == nil {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code: "401",
				Message: "Unauthorized",
			})
	}

	user, err := ur.GetByUserID(token.UserID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	if user == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code: "404",
				Message: "User Not Found",
			})
	}

	// 同性は見れない
	oppositeUser, err := ur.GetByUserID(p.UserID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	if oppositeUser == nil {
		return si.NewGetProfileByUserIDBadRequest().WithPayload(
			&si.GetProfileByUserIDBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}


	if user.Gender == oppositeUser.Gender {
		return si.NewGetProfileByUserIDBadRequest().WithPayload(
			&si.GetProfileByUserIDBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	
	sEnt := user.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

// update
func PutProfile(p si.PutProfileParams) middleware.Responder {
	tr := repositories.NewUserTokenRepository()
	ur := repositories.NewUserRepository()

	token, err := tr.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error (in get token)",
			})
	}

	if token == nil {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code: "401",
				Message: "Unauthorized",
			})
	}

	if token.UserID != p.UserID {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code: "403",
				Message: "Forbidden",
			})
	}

	usr, err := ur.GetByUserID(token.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error (in get user)",
			})
	}

	if usr == nil {
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code: "400",
				Message: "Bad Request (in get user)",
			})
	}

	usr.ApplyParams(p.Params)
	err = ur.Update(usr)

	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
			Code: "500",
			Message: "Internal Server Error",
			})
	}

	updatedUser, err := ur.GetByUserID(token.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
			Code: "500",
			Message: "Internal Server Error",
			})
	}
	
	sUsr := updatedUser.Build()

	return si.NewPutProfileOK().WithPayload(&sUsr)
}
