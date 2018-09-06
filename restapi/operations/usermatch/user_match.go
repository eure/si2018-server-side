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
	limit := int(p.Limit)
	offset := int(p.Offset)
	if limit < 0 || offset < 0 {
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				"400",
				"Bad Request",
			})
	}
	entMatch, errMatch := rMatch.FindByUserIDWithLimitOffset(sEntToken.UserID, limit, offset)
	if errMatch != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	matches := entities.UserMatches(entMatch)
	sMatches := matches.Build()

	rUser := repositories.NewUserRepository()

	// 上で取得した全てのpartnerIDについて、プロフィール情報を取得してpayloadsに格納する。

	var IDs []int64
	for _, sMatch := range sMatches{
		IDs = append(IDs, sMatch.PartnerID)
	}

	partners, errFind := rUser.FindByIDs(IDs)

	if errFind != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	entPartners := entities.Users(partners)
	sEntPartners := entPartners.Build() // プロフィールのリスト

	var payloads []*models.MatchUserResponse
	for _, sEntPartner := range sEntPartners{
		//entities.User -> models.LikeUserResponse
		r := models.MatchUserResponse{}
		r.ID = sEntPartner.ID
		r.Nickname = sEntPartner.Nickname
		r.Tweet = sEntPartner.Tweet
		r.Introduction = sEntPartner.Introduction
		r.ResidenceState = sEntPartner.ResidenceState
		r.HomeState = sEntPartner.HomeState
		r.Education = sEntPartner.Education
		r.Job = sEntPartner.Job
		r.AnnualIncome = sEntPartner.AnnualIncome
		r.Height = sEntPartner.Height
		r.BodyBuild = sEntPartner.BodyBuild
		r.MaritalStatus = sEntPartner.MaritalStatus
		r.Child = sEntPartner.Child
		r.WhenMarry = sEntPartner.WhenMarry
		r.WantChild = sEntPartner.WantChild
		r.Smoking = sEntPartner.Smoking
		r.Drinking = sEntPartner.Drinking
		r.Holiday = sEntPartner.Holiday
		r.HowToMeet = sEntPartner.HowToMeet
		r.CostOfDate = sEntPartner.CostOfDate
		r.NthChild = sEntPartner.NthChild
		r.Housework = sEntPartner.Housework
		r.ImageURI = sEntPartner.ImageURI
		r.CreatedAt = sEntPartner.CreatedAt
		r.UpdatedAt = sEntPartner.UpdatedAt
		/* r.LikedAt = (探しても見つからない)*/
		payloads = append(payloads,&r)
	}

	return si.NewGetMatchesOK().WithPayload(payloads)
}
