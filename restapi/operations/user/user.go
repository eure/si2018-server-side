package user

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	nutr := repositories.NewUserTokenRepository()
	nulr := repositories.NewUserLikeRepository()
	nur := repositories.NewUserRepository()
	//ngur := repositories.NewGetUserRepository()
	// find userid
	usrid, err := nutr.GetByToken(p.Token)
	if err != nil {
		return GetUserRespUnauthErr()
	}
	// find userlike
	userlike, err := nulr.FindLikeAll(usrid.UserID)
	if err != nil {
		return GetUserRespInternalErr()
	}
	// find user
	userdesc, err := nur.GetByUserID(usrid.UserID)
	if err != nil {
		return GetUserRespInternalErr()
	}
	ent, err := nur.FindWithCondition(int(p.Limit), int(p.Offset), userdesc.GetOppositeGender(), userlike)
	if err != nil {
		return GetUserRespInternalErr()
	}

	var ud []*models.User
	for i := 0; i < len(ent); i++ {
		us := ent[i].Build()
		ud = append(ud, &us)
	}

	return si.NewGetUsersOK().WithPayload(ud)

}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	nur := repositories.NewUserRepository()
	nut := repositories.NewUserTokenRepository()
	usrid, err := nur.GetByUserID(p.UserID)
	if err != nil {
		return GetProfileInternalErr()
	}
	usrtoken, err := nut.GetByToken(p.Token)
	if err != nil {
		return GetProfileRespUnauthErr()

	}
	if usrid.ID != usrtoken.UserID {
		return GetProfileNotFoundErr()
	}
	sEnt := usrid.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	nur := repositories.NewUserRepository()
	//Find user
	user, err := nur.GetByUserID(p.UserID)
	if err != nil {
		return PutProfileInternalErr()
	}
	binduser(p.Params, user)
	err = nur.Update(user)
	if err != nil {
		return PutProfileForbiddenErr()
	}
	// Want Response User profile
	respuser, err := nur.GetByUserID(p.UserID)
	if err != nil {
		return PutProfileInternalErr()
	}
	updateuser := respuser.Build()
	return si.NewPutProfileOK().WithPayload(&updateuser)
}

func binduser(user si.PutProfileBody, ent *entities.User) {
	ent.AnnualIncome = user.AnnualIncome
	ent.BodyBuild = user.BodyBuild
	ent.Child = user.Child
	ent.CostOfDate = user.CostOfDate
	ent.Drinking = user.Drinking
	ent.Education = user.Education
	ent.Height = user.Height
	ent.Holiday = user.Holiday
	ent.HomeState = user.HomeState
	ent.Housework = user.Housework
	ent.HowToMeet = user.HowToMeet
	ent.ImageURI = user.ImageURI
	ent.Introduction = user.Introduction
	ent.Job = user.Job
	ent.MaritalStatus = user.MaritalStatus
	ent.Nickname = user.Nickname
	ent.NthChild = user.NthChild
	ent.ResidenceState = user.ResidenceState
	ent.Smoking = user.Smoking
	ent.Tweet = user.Tweet
	ent.WantChild = user.WantChild
	ent.WhenMarry = user.WhenMarry

}

func GetUserRespUnauthErr() middleware.Responder {
	return si.NewGetUsersUnauthorized().WithPayload(
		&si.GetUsersUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

func GetUserRespBadReqestErr() middleware.Responder {
	return si.NewGetUsersBadRequest().WithPayload(
		&si.GetUsersBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

func GetUserRespInternalErr() middleware.Responder {
	return si.NewGetUsersInternalServerError().WithPayload(
		&si.GetUsersInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func PutProfileRespUnauthErr() middleware.Responder {
	return si.NewPutProfileUnauthorized().WithPayload(
		&si.PutProfileUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

func PutProfileBadRequestErr() middleware.Responder {
	return si.NewPutProfileBadRequest().WithPayload(
		&si.PutProfileBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

func PutProfileForbiddenErr() middleware.Responder {
	return si.NewPutProfileForbidden().WithPayload(
		&si.PutProfileForbiddenBody{
			Code:    "403",
			Message: "Forbidden",
		})
}

func PutProfileInternalErr() middleware.Responder {
	return si.NewPutProfileInternalServerError().WithPayload(
		&si.PutProfileInternalServerErrorBody{
			Code:    "500",
			Message: "Intsernal Server Error",
		})
}

func GetProfileRespUnauthErr() middleware.Responder {
	return si.NewGetProfileByUserIDUnauthorized().WithPayload(
		&si.GetProfileByUserIDUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

func GetProfileBadRequestErr() middleware.Responder {
	return si.NewGetProfileByUserIDBadRequest().WithPayload(
		&si.GetProfileByUserIDBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

func GetProfileNotFoundErr() middleware.Responder {
	return si.NewGetProfileByUserIDNotFound().WithPayload(
		&si.GetProfileByUserIDNotFoundBody{
			Code:    "404",
			Message: "User Not Found",
		})
}

func GetProfileInternalErr() middleware.Responder {
	return si.NewGetProfileByUserIDInternalServerError().WithPayload(
		&si.GetProfileByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}
