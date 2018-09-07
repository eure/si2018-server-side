package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	/*
		1. tokenのvalidation
		2. tokenからuseridを取得
		3. useridからマッチングしたユーザーの一覧を取得
		// userIDはいいねを送った人, partnerIDはいいねを受け取った人
	*/

	// Tokenがあるかどうか
	if p.Token == "" {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Required",
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
				Message: "Token Is Invalid",
			})
	}

	sEntToken := entToken.Build()

	//useridからマッチングしたユーザーの一覧を取得
	rMatch := repositories.NewUserMatchRepository()
	limit := int(p.Limit)
	offset := int(p.Offset)
	if limit <= 0 || offset < 0 {
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

	partnerMatchedAt := map[int64]strfmt.DateTime{}

	for _, sMatch := range sMatches {
		partnerMatchedAt[sMatch.PartnerID] = sMatch.CreatedAt
	}

	// 上で取得した全てのpartnerIDについて、プロフィール情報を取得してpayloadsに格納する。

	var IDs []int64
	for _, sMatch := range sMatches {
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

	var payloads []*models.MatchUserResponse
	for _, partner := range partners {
		var r entities.MatchUserResponse
		r.ApplyUser(partner)
		r.MatchedAt = partnerMatchedAt[partner.ID]
		m := r.Build()
		payloads = append(payloads, &m)
	}

	return si.NewGetMatchesOK().WithPayload(payloads)
}
