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
				Message: "Internal Server Error: GetByToken failed: " + err.Error(),
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
				Message: "Internal Server Error: GetByUserID failed: " + err.Error(),
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
				Message: "Internal Server Error: FindLikeAll failed: " + err.Error(),
			})
	}
	for _, id := range like {
		idmap[id] = true
	}
	matched, err := rm.FindAllByUserID(u.ID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: FindAllByUserID failed: " + err.Error(),
			})
	}
	for _, id := range matched {
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
				Message: "Internal Server Error: FindWithCondition failed: " + err.Error(),
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
				Message: "Internal Server Error: GetByUserID failed: " + err.Error(),
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
				Message: "Internal Server Error: GetByToken failed: " + err.Error(),
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
				Message: "Internal Server Error: GetByUserID failed: " + err.Error(),
			})
	}
	if u == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request: GetByUserID",
			})
	}
	sEnt := u.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	r := repositories.NewUserRepository()
	rt := repositories.NewUserTokenRepository()
	token, err := rt.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByUserID failed: " + err.Error(),
			})
	}
	if token == nil {
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code:    "400",
				Message: "Bad Request: GetByUserID failed",
			})
	}

	t, err := rt.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByToken failed: " + err.Error(),
			})
	}
	if t == nil {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗): GetByToken failed",
			})
	}
	if t.UserID != p.UserID {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code:    "403",
				Message: "Forbidden. (他の人のプロフィールは更新できません.): Token does not match",
			})
	}
	u, err := r.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: Get By UserID failed: " + err.Error(),
			})
	}
	if u == nil {
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code:    "400",
				Message: "Bad Request: GetByUserID failed",
			})
	}
	params := p.Params
	if u.MaritalStatus != "独身(未婚)" && params.MaritalStatus == "独身(未婚)" {
		// 未婚でない人が未婚になることはありえないはず
		// システム上、結婚歴で嘘をつくことを認めるかは議論の余地あり
		// とりあえず今回は見逃す
	}
	u.AnnualIncome = params.AnnualIncome
	u.BodyBuild = params.BodyBuild
	u.Child = params.Child
	u.CostOfDate = params.CostOfDate
	u.Drinking = params.Drinking
	u.Education = params.Education
	u.Height = params.Height
	u.Holiday = params.Holiday
	u.HomeState = params.HomeState
	u.Housework = params.Housework
	u.HowToMeet = params.HowToMeet
	u.ImageURI = params.ImageURI
	u.Introduction = params.Introduction
	u.Job = params.Job
	u.MaritalStatus = params.MaritalStatus
	u.Nickname = params.Nickname
	u.NthChild = params.NthChild
	u.ResidenceState = params.ResidenceState
	u.Smoking = params.Smoking
	u.Tweet = params.Tweet
	u.WantChild = params.WantChild
	u.WhenMarry = params.WhenMarry
	err = r.Update(u)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: Update failed: " + err.Error(),
			})
	}
	u, err = r.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: Get By UserID failed: " + err.Error(),
			})
	}
	if u == nil {
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code:    "400",
				Message: "Bad Request: GetByUserID failed",
			})
	}
	sEnt := u.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}
