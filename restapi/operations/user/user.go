package user

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	r := repositories.NewUserRepository()

	user, _ := r.GetByToken(p.Token)

	liked := repositories.NewUserLikeRepository()
	users, _ := liked.FindLikeAll(user.ID)

	var ent entities.Users
	ent, _ = r.FindWithCondition(int(p.Limit),int(p.Offset),user.GetOppositeGender(),users)

	sEnt := ent.Build()
	return si.NewGetUsersOK().WithPayload(sEnt)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	r := repositories.NewUserRepository()

	ent, err := r.GetByUserID(p.UserID)

	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if ent == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	r := repositories.NewUserRepository()
	t := repositories.NewUserTokenRepository()

	token , _ := t.GetByUserID(p.UserID)
	user , _ := r.GetByUserID(p.UserID)
	if token.UserID == p.UserID {
		ProfileUpdate(user,p.Params)

		// 書き換えた情報でデータベースを更新
		err := r.Update(user)
		if err != nil {
			return si.NewPutProfileInternalServerError().WithPayload(
				&si.PutProfileInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	} else {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code:    "403",
				Message: "Forbidden",
			})
	}

	ent, _ := r.GetByUserID(p.UserID)
	sEnt := ent.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}


func ProfileUpdate(u *entities.User, p si.PutProfileBody) {
	// annual income
	u.AnnualIncome = p.AnnualIncome
	// body build
	u.BodyBuild = p.BodyBuild
	// child
	u.Child = p.Child
	// cost of date
	u.CostOfDate = p.CostOfDate
	// drinking
	u.Drinking = p.Drinking
	// education
	u.Education = p.Education
	// height
	u.Height = p.Height
	// holiday
	u.Holiday = p.Holiday
	// home state
	u.HomeState = p.HomeState
	// housework
	u.Housework = p.Housework
	// how to meet
	u.HowToMeet = p.HowToMeet
	// image uri
	u.ImageURI = p.ImageURI
	// introduction
	u.Introduction = p.Introduction
	// job
	u.Job = p.Job
	// marital status
	u.MaritalStatus = p.MaritalStatus
	// nickname
	u.Nickname = p.Nickname
	// nth child
	u.NthChild = p.NthChild
	// residence state
	u.ResidenceState = p.ResidenceState
	// smoking
	u.Smoking = p.Smoking
	// tweet
	u.Tweet = p.Tweet
	// want child
	u.WantChild = p.WantChild
	// when marry
	u.WhenMarry = p.WhenMarry
}