package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	userTokenEnt, err := repositories.NewUserTokenRepository().GetByToken(p.Token)

	if err != nil {
		return getMatchesInternalServerErrorRespose()
	}

	if userTokenEnt == nil {
		return getMatchesUnauthorizedResponse()
	}

	// int64になっているのでcastする必要がある
	limit := int(p.Limit)
	offset := int(p.Offset)
	userID := userTokenEnt.UserID

	userMatchRepository := repositories.NewUserMatchRepository()
	userMatches, err := userMatchRepository.FindByUserIDWithLimitOffset(userID, limit, offset)

	if err != nil {
		return getMatchesInternalServerErrorRespose()
	}

	var matchUserResponsesEnt entities.MatchUserResponses

	for _, userMatchEnt := range userMatches {
		matchUserResponse := entities.MatchUserResponse{MatchedAt: userMatchEnt.CreatedAt}

		//FIXME ループ内でクエリ発行は最低の行為のような気がする
		user, err := repositories.NewUserRepository().GetByUserID(userMatchEnt.UserID)
		if err != nil {
			return getMatchesInternalServerErrorRespose()

		}

		matchUserResponse.ApplyUser(*user)

		matchUserResponsesEnt = append(matchUserResponsesEnt, matchUserResponse)
	}

	mathUserResponses := matchUserResponsesEnt.Build()

	return si.NewGetMatchesOK().WithPayload(mathUserResponses)
}

func getMatchesInternalServerErrorRespose() middleware.Responder {
	return si.NewGetMatchesInternalServerError().WithPayload(
		&si.GetMatchesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getMatchesUnauthorizedResponse() middleware.Responder {
	return si.NewGetMatchesUnauthorized().WithPayload(
		&si.GetMatchesUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}
