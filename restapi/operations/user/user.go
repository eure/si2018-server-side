package user

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	u := repositories.NewUserRepository()
	l := repositories.NewUserLikeRepository()
	t := repositories.NewUserTokenRepository()

	// tokenから UserToken entitiesを取得 (Validation)
	token := p.Token
	loginUserToken , err := t.GetByToken(token)
	if err != nil {
		return outPutGetStatus(500)
	}
	if loginUserToken == nil {
		return outPutGetStatus(401)
	}
	
	// limit が20かどうか検出
	if p.Limit != int64(20) {
		return outPutGetStatus(400)
	}
	// offset が0以上かどうか検出
	if p.Offset < int64(0) {
		return outPutGetStatus(400)
	}
	
	loginUser, err := u.GetByUserID(loginUserToken.UserID)
	if err != nil {
		return outPutGetStatus(500)
	}
	if loginUser == nil {
		return outPutGetStatus(400)
	}
	
	// すでにいいね！しているユーザーのUserID int64を集める
	likedUserIDs, err := l.FindLikePart(loginUserToken.UserID)

	if err != nil {
		return outPutGetStatus(500)
	}
	if likedUserIDs == nil {
		return outPutGetStatus(400)
	}
	
	var ent entities.Users
	ent, err = u.FindWithCondition(int(p.Limit),int(p.Offset),loginUser.GetOppositeGender(),likedUserIDs)
	
	if err != nil {
		return outPutGetStatus(500)
	}

	if ent == nil {
		return outPutGetStatus(400)
	}
	
	sEnt := ent.Build()
	return si.NewGetUsersOK().WithPayload(sEnt)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	r := repositories.NewUserRepository()
	t := repositories.NewUserTokenRepository()
	
	// tokenから UserToken entitiesを取得 (Validation)
	token := p.Token
	loginUser , err := t.GetByToken(token)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if loginUser == nil {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗)",
			})
	}
	
	// Paramsから Profileを取得したいユーザーのIDを取得
	profileUserID := p.UserID
	ent, err := r.GetByUserID(profileUserID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if ent == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found",
			})
	}

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	r := repositories.NewUserRepository()
	t := repositories.NewUserTokenRepository()

	// tokenから UserToken entitiesを取得 (Validation)
	token := p.Params.Token
	loginUserToken , err := t.GetByToken(token)
	if err != nil {
		return outPutPutStatus(401)
	}
	if loginUserToken == nil {
		return outPutPutStatus(400)
	}
	
	loginUser , err := r.GetByUserID(loginUserToken.UserID)
	
	if err != nil {
		return outPutPutStatus(500)
	}
	if loginUserToken.UserID == p.UserID {
		ProfileUpdate(loginUser,p.Params)

		// 書き換えた情報でデータベースを更新
		err := r.Update(loginUser)
		if err != nil {
			return outPutPutStatus(500)
		}
	} else {
		return outPutPutStatus(403)
	}

	ent, err := r.GetByUserID(p.UserID)
	if err != nil {
		return outPutPutStatus(500)
	}
	if ent == nil {
		return outPutPutStatus(400)
	}
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

func outPutGetStatus (num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	case 401:
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗)",
			})
	case 400:
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	return nil
}

func outPutPutStatus (num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	case 403:
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code:    "403",
				Message: "Forbidden",
			})
	case 401:
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗)",
			})
	case 400:
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	return nil
}