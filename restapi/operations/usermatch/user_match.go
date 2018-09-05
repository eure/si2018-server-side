package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	/*
	1. tokenのvalidation
	2. tokenからuseridを取得
	3. useridからマッチングしたユーザーの一覧を取得
	*/


	// Tokenがあるかどうか
	if p.Token == "" {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code:    "401",
				Message: "No Token",
			})
	}

	// tokenからuserIDを取得
	rToken := repositories.NewUserTokenRepository()
	entToken, errToken := rToken.GetByToken(p.Token)
	if errToken != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entToken == nil {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized Token",
			})
	}

	sEntToken := entToken.Build()
	//useridからマッチングしたユーザーの一覧を取得
	rMatch := repositories.NewUserMatchRepository()
	entMatch, errMatch := rMatch.FindByUserIDWithLimitOffset(sEntToken.UserID, int(p.Limit),int(p.Offset))
	if errMatch != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	matches := entities.UserMatches(entMatch)
	sMatches := matches.Build()

	// sUsersが全て同じになってしまう。ポインタについては解決方法がわからない。

	rUser := repositories.NewUserRepository()

	// 上で取得した全てのpartnerIDについて、プロフィール情報を取得してpayloadsに格納する。
	var payloads []*models.MatchUserResponse
	for _, sMatch := range sMatches{
		has, err := rUser.GetByUserID(sMatch.PartnerID)

		if err != nil{
			return si.NewGetMatchesInternalServerError().WithPayload(
				&si.GetMatchesInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
		//entities.User -> models.LikeUserResponse
		r := models.MatchUserResponse{}
		r.ID = has.ID
		r.Nickname = has.Nickname
		r.Tweet = has.Tweet
		r.Introduction = has.Introduction
		r.ResidenceState = has.ResidenceState
		r.HomeState = has.HomeState
		r.Education = has.Education
		r.Job = has.Job
		r.AnnualIncome = has.AnnualIncome
		r.Height = has.Height
		r.BodyBuild = has.BodyBuild
		r.MaritalStatus = has.MaritalStatus
		r.Child = has.Child
		r.WhenMarry = has.WhenMarry
		r.WantChild = has.WantChild
		r.Smoking = has.Smoking
		r.Drinking = has.Drinking
		r.Holiday = has.Holiday
		r.HowToMeet = has.HowToMeet
		r.CostOfDate = has.CostOfDate
		r.NthChild = has.NthChild
		r.Housework = has.Housework
		r.ImageURI = has.ImageURI
		r.CreatedAt = has.CreatedAt
		r.UpdatedAt = has.UpdatedAt
		/* r.LikedAt = (探しても見つからない)*/
		payloads = append(payloads,&r)
	}
	return si.NewGetMatchesOK().WithPayload(payloads)
}
