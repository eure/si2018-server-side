package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	//FIXME Token認証は共通化したい
	r := repositories.NewUserTokenRepository()

	ent, err := r.GetByToken(p.Token)

	if err != nil {
		return GetUsersInteralServerErrorResponse("Internal Server Error")
	}
	if ent == nil {
		return GetUsersUnauthorizedResponse("Your Token Is Invalid")
	}

	userToken := ent.Build()
	userID := userToken.UserID

	userLikeRepository := repositories.NewUserLikeRepository()
	exclusionIds, err := userLikeRepository.FindLikeAll(userID)

	if err != nil {
		return GetUsersInteralServerErrorResponse("Internal Server Error")
	}

	userRepository := repositories.NewUserRepository()
	userEnt, err := userRepository.GetByUserID(userID)
	if err != nil {
		return GetUsersInteralServerErrorResponse("Internal Server Error")
	}

	// int64になっているのでcastする必要がある
	limit := int(p.Limit)
	offset := int(p.Offset)
	gender := userEnt.GetOppositeGender()

	var usersEnt entities.Users
	usersEnt, err = userRepository.FindUsers(limit, offset, gender, exclusionIds)
	if err != nil {
		return GetUsersInteralServerErrorResponse("Internal Server Error")
	}

	var ids []int64
	for _, u := range usersEnt {
		ids = append(ids, u.ID)
	}

	userImageRepository := repositories.NewUserImageRepository()
	images, err := userImageRepository.GetByUserIDs(ids)
	if err != nil {
		return GetUsersInteralServerErrorResponse("Internal Server Error")
	}

	for i := range usersEnt {
		usersEnt[i].ImageURI = images[i].Path
	}

	users := usersEnt.Build()
	return si.NewGetUsersOK().WithPayload(users)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	//FIXME Token認証は共通化したい
	r := repositories.NewUserTokenRepository()

	userTokenEnt, err := r.GetByToken(p.Token)

	if err != nil {
		return GetUserProfileByUserIDInternalServerErrorResponse("Internal Server Error")
	}
	if userTokenEnt == nil {
		return GetUsersUnauthorizedResponse("Your Token Is Invalid")
	}

	userID := p.UserID

	userEnt, err := repositories.NewUserRepository().GetByUserID(userID)
	if err != nil {
		return GetUserProfileByUserIDInternalServerErrorResponse("Internal Server Error")
	}

	if userEnt == nil {
		return GetUserProfileByUserIDNotFoundResponse("User Not Found")
	}

	userImageRepository := repositories.NewUserImageRepository()
	userImage, err := userImageRepository.GetByUserID(userID)
	if err != nil {
		return GetUserProfileByUserIDInternalServerErrorResponse("Internal Server Error")
	}

	userEnt.ImageURI = userImage.Path

	user := userEnt.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&user)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	userTokenRepository := repositories.NewUserTokenRepository()
	userTokenEnt, err := userTokenRepository.GetByToken(p.Params.Token)

	if err != nil {
		return PutProfileInternalServerErrorResponse()
	}

	if userTokenEnt == nil {
		return PutProfileUnauthorizedResponse()
	}

	userID := p.UserID

	if userTokenEnt.UserID != userID {
		return PutProfileForbiddenResponse()
	}

	userEnt, err := repositories.NewUserRepository().GetByUserID(userID)

	if err != nil || userEnt == nil {
		return PutProfileInternalServerErrorResponse()
	}

	updateUserEnt := repositories.NewUserRepository().ParamsToUserEnt(userEnt, p.Params)

	err = repositories.NewUserRepository().Update(updateUserEnt)
	if err != nil {
		return PutProfileInternalServerErrorResponse()
	}

	updatedUserEnt, err := repositories.NewUserRepository().GetByUserID(userID)
	if err != nil {
		return PutProfileInternalServerErrorResponse()
	}

	user := updatedUserEnt.Build()
	return si.NewPutProfileOK().WithPayload(&user)
}
