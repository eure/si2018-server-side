package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/restapi/operations/util"
	"log"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	limit := p.Limit
	offset := p.Offset
	token := p.Token

	// Validations
	err1 := util.ValidateLimit(limit)
	err2 := util.ValidateOffset(offset)
	if (err1 != nil) || (err2 != nil) {
		log.Print("Limit/Offset validation error")
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
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

	id, _ := util.GetIDByToken(token)
	user, err := ru.GetByUserID(id)

	likes, err := rl.FindLikeAll(id) /* TODO filter? */
	if err != nil {
		log.Print("Get all likes error")
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	users_ent, err := ru.FindWithCondition(int(limit), int(offset), user.GetOppositeGender(), likes);
	if err != nil {
		log.Print("Find user error", err)
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
	
	// Validation
	err := util.ValidateToken(token)
	if err != nil {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	myid, _ := util.GetIDByToken(token)

	r := repositories.NewUserRepository()

	me, err1 := r.GetByUserID(myid)
	user, err2 := r.GetByUserID(id)
	if err1 != nil || err2 != nil {
		log.Print("Get user error")
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	// User exists?
	if user == nil { 
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found",
			})
	}
	// Same gender?
	if user.Gender == me.Gender && myid != id {
		log.Print("Same gender")
		return si.NewGetProfileByUserIDBadRequest().WithPayload(
			&si.GetProfileByUserIDBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	
	// Prepare response
	res := user.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&res)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	id := p.UserID
	ps := p.Params
	token := ps.Token
	
	// Validation
	err := util.ValidateToken(token)
	if err != nil {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// Same id in token and param?
	tid, _ := util.GetIDByToken(token)
	if id != tid {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code:    "403",
				Message: "Forbidden",
			})
	}

	// Prepare new User to update
	r := repositories.NewUserRepository()
	up, _ := r.GetByUserID(id)
	user := *up

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

	// Update
	err = r.Update(&user)
	if err != nil {
		log.Print("Update err")
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// Response (return to check wheter updated or not)
	after, _ := r.GetByUserID(id)
	res := after.Build()

	return si.NewPutProfileOK().WithPayload(&res)
}