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
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByToken failed",
			})
	}
	if t == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗): GetByToken failed",
			})
	}
	u, err := r.GetByUserID(t.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByUserID failed",
			})
	}
	if u == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request: GetByUserID failed",
			})
	}
	idmap := make(map[int64]bool)
	like, err := rl.FindLikeAll(u.ID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: FindLikeAll failed",
			})
	}
	if like == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request: FindLikeAll failed",
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
				Message: "Internal Server Error: FindAllByUserID failed",
			})
	}
	if mached == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request: FindAllByUserID failed",
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
	if ent == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request: FindWithCondition failed",
			})
	}
	hoge := entities.Users(ent)

	return si.NewGetUsersOK().WithPayload(hoge.Build())
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	r := repositories.NewUserRepository()
	rt := repositories.NewUserTokenRepository()

	token, err := rt.GetByUserID(p.UserID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByUserID failed",
			})
	}
	if token == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found. (そのIDのユーザーは存在しません.): GetByUserID failed",
			})
	}

	t, err := rt.GetByToken(p.Token)
	if err != nil {
		return si.NewGetTokenByUserIDInternalServerError().WithPayload(
			&si.GetTokenByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByToken failed",
			})
	}
	if t == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗): GetByToken failed",
			})
	}
	if t.UserID != p.UserID {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗): Token does not match",
			})
	}
	u, err := r.GetByUserID(p.UserID)
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
			Message: "Bad Request: GetByUserID",
		})
	}
	sEnt := u.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	return si.NewPutProfileOK()
}
