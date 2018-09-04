package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/restapi/operations/util"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	limit := p.Limit
	offset := p.Offset
	token := p.Token
	/* TODO bad request? */
	
	err := util.ValidateToken(token)
	if err != nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Your Token Is Invalid",
			})
	}
	
	ru := repositories.NewUserRepository()
	rl := repositories.NewUserLikeRepository()

	id, err := util.GetIDByToken(token)
	user, err := ru.GetByUserID(id)
	likes, err := rl.FindLikeAll(id)

	users_ent, err := ru.FindWithCondition(int(limit), int(offset), user.GetOppositeGender(), likes);
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	var users entities.Users
	users = entities.Users(users_ent)
	return si.NewGetUsersOK().WithPayload(users.Build())
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	id := p.UserID
	token := p.Token
	/* TODO bad request? */
	
	err := util.ValidateToken(token)
	if err != nil {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	r := repositories.NewUserRepository()

	user, err := r.GetByUserID(id)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if user == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found",
			})
	}
	
	res := user.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&res)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	id := p.UserID
	ps := p.Params
	token := ps.Token
	
	err := util.ValidateToken(token)
	if err != nil {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	tid, _ := util.GetIDByToken(token)
	if id != tid {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code:    "403",
				Message: "Forbidden",
			})
	}

	r := repositories.NewUserRepository()
	up, _ := r.GetByUserID(id)
	user := *up

	//fmt.Printf("user:%+v\n", user)
	//fmt.Printf("param:%+v\n", p.Params)

	user.AnnualIncome = ps.AnnualIncome
	user.BodyBuild = ps.BodyBuild
	user.Child = ps.Child
	user.CostOfDate = ps.CostOfDate
	user.Drinking = ps.Drinking
	user.Education = ps.Education
	user.Height = ps.Height
	user.Holiday = ps.Holiday
	user.HomeState = ps.HomeState
	user.Housework = ps.Housework
	user.HowToMeet = ps.HowToMeet
	user.ImageURI = ps.ImageURI
	user.Introduction = ps.Introduction
	user.Job = ps.Job
	user.MaritalStatus = ps.MaritalStatus
	user.Nickname = ps.Nickname
	user.NthChild = ps.NthChild
	user.ResidenceState = ps.ResidenceState
	user.Smoking = ps.Smoking
	user.Tweet = ps.Tweet
	user.WantChild = ps.WantChild
	user.WhenMarry = ps.WhenMarry

	//fmt.Printf("newu:%+v\n", user)
	err = r.Update(&user)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	after, _ := r.GetByUserID(id)
	res := after.Build()

	return si.NewPutProfileOK().WithPayload(&res)
}