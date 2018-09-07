package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"sort"
	"time"
)

type UserResponses []*models.MatchUserResponse

func (a UserResponses) Len() int      { return len(a) }
func (a UserResponses) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a UserResponses) Less(i, j int) bool {
	ai := time.Time(a[i].MatchedAt)
	aj := time.Time(a[j].MatchedAt)
	return !ai.Before(aj)
}

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
				Code:    "400",
				Message: "Bad Request",
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

	var IDs []int64
	partnerMatchedAt := map[int64]strfmt.DateTime{}
	// id -- 時間の対応mapとpartneridのリストの作成
	for _, sMatch := range sMatches {
		partnerMatchedAt[sMatch.PartnerID] = sMatch.CreatedAt
		IDs = append(IDs, sMatch.PartnerID)
	}

	if IDs == nil {
		var payloads []*models.MatchUserResponse
		return si.NewGetMatchesOK().WithPayload(payloads)
	}

	// 上で取得した全てのpartnerIDについて、プロフィール情報と画像URIを取得してpayloadsに格納する。
	partners, errFind := rUser.FindByIDs(IDs)

	if errFind != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 画像URI取得
	rImage := repositories.NewUserImageRepository()
	entImages, errImages := rImage.GetByUserIDs(IDs)
	if errImages != nil || entImages == nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// id -- pathの対応リストを作成
	idPaths := map[int64]string{}
	for _, entImage := range entImages {
		idPaths[entImage.UserID] = entImage.Path
	}

	var payloads []*models.MatchUserResponse
	for _, partner := range partners {
		var r entities.MatchUserResponse
		r.ApplyUser(partner)
		r.MatchedAt = partnerMatchedAt[partner.ID]
		r.ImageURI = idPaths[partner.ID]
		m := r.Build()
		payloads = append(payloads, &m)
	}

	sort.Sort(UserResponses(payloads))

	return si.NewGetMatchesOK().WithPayload(payloads)
}
