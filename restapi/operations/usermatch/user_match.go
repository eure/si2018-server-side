package usermatch

import (
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"

	"fmt"

	"github.com/eure/si2018-server-side/models"
	"github.com/go-openapi/strfmt"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	tokenRepo := repositories.NewUserTokenRepository()
	matchRepo := repositories.NewUserMatchRepository()
	userRepo := repositories.NewUserRepository()

	// tokenが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if tokenOwner == nil {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// トークンの持ち主がマッチングしているお相手の一覧を取得します。
	matches, ids, err := matchRepo.FindByUserIDWithLimitOffset(tokenOwner.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	var matchedAt []strfmt.DateTime
	for _, match := range matches {
		matchedAt = append(matchedAt, match.CreatedAt)
	}

	// models.MatchUserResponse と entities.User をマッピングします。
	matchUsers, err := userRepo.FindByIDs(ids)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	var matchUserResponses []*models.MatchUserResponse
	for i, matchUser := range matchUsers {
		matchUserResponses = append(matchUserResponses, &models.MatchUserResponse{
			MatchedAt:      matchedAt[i],
			AnnualIncome:   matchUser.AnnualIncome,
			Birthday:       matchUser.Birthday,
			BodyBuild:      matchUser.BodyBuild,
			Child:          matchUser.Child,
			CostOfDate:     matchUser.CostOfDate,
			Drinking:       matchUser.Drinking,
			Education:      matchUser.Education,
			Gender:         matchUser.Gender,
			Height:         matchUser.Height,
			Holiday:        matchUser.Holiday,
			HomeState:      matchUser.HomeState,
			Housework:      matchUser.Housework,
			HowToMeet:      matchUser.HowToMeet,
			ID:             matchUser.ID,
			ImageURI:       matchUser.ImageURI,
			Introduction:   matchUser.Introduction,
			Job:            matchUser.Job,
			MaritalStatus:  matchUser.MaritalStatus,
			Nickname:       matchUser.Nickname,
			NthChild:       matchUser.NthChild,
			ResidenceState: matchUser.ResidenceState,
			Smoking:        matchUser.Smoking,
			Tweet:          matchUser.Tweet,
			WantChild:      matchUser.WantChild,
			WhenMarry:      matchUser.WhenMarry,
			CreatedAt:      matchUser.CreatedAt,
			UpdatedAt:      matchUser.UpdatedAt,
		})
		fmt.Println(matchUser.ID, matchedAt[i])
	}

	return si.NewGetMatchesOK().WithPayload(matchUserResponses)
}
