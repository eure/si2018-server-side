package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	rt := repositories.NewUserTokenRepository()
	r := repositories.NewUserRepository()
	rl := repositories.NewUserLikeRepository()
	rm := repositories.NewUserMatchRepository()
	t, err := rt.GetByToken(p.Token)
	if err != nil {
		return si.NewGetTokenByUserIDInternalServerError().WithPayload(
			&si.GetTokenByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: Get By Token failed",
			})
	}
	if t == nil {
		return si.NewGetUsersUnauthorized().WithPayload(&si.GetUsersUnauthorizedBody{
			Code:    "401",
			Message: "Bad Request: Invalid Token",
		})
	}
	u, err := r.GetByUserID(t.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: Get By UserID failed",
			})
	}
	if u == nil {
		return si.NewGetUsersBadRequest().WithPayload(&si.GetUsersBadRequestBody{
			Code:    "400",
			Message: "Bad Request: User Not Found",
		})
	}
	idmap := make(map[int64]bool)
	like, err := rl.FindLikeAll(u.ID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: Find All Likes failed",
			})
	}
	for _, id := range like {
		idmap[id] = true
	}
	mached, err := rm.FindAllByUserID(u.ID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: Find All Matches failed",
			})
	}
	for _, id := range mached {
		idmap[id] = true
	}
	ids := make([]int64, 0)
	for k := range idmap {
		ids = append(ids, k)
	}
	ent, err := r.FindWithCondition(int(p.Limit), int(p.Offset), u.GetOppositeGender(), ids)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: Find With Condition failed",
			})
	}
	hoge := entities.Users(ent)

	return si.NewGetUsersOK().WithPayload(hoge.Build())
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	return si.NewGetProfileByUserIDOK()
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	return si.NewPutProfileOK()
}
