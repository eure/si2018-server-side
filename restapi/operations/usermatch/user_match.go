package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	usertokenHandler := repositories.NewUserTokenRepository()
	usermatchHandler := repositories.NewUserMatchRepository()
	userHandler := repositories.NewUserRepository()

	token := p.Token
	limit := int(p.Limit)
	offset := int(p.Offset)
	// Find user using token
	user, err := usertokenHandler.GetByToken(token)
	if err != nil {
		return GetMatchessRespUnauthErr()
	}
	// Does exist user?
	if user == nil {
		return GetMatchessRespBadReqestErr()
	}

	// Find all match user
	ent, err := usermatchHandler.FindByUserIDWithLimitOffset(user.UserID, limit, offset)
	if err != nil {
		return GetMatchesRespInternalErr()
	}

	// If not matcheing to anyone
	if len(ent) == 0 {
		var notmatches entities.MatchUserResponses
		resp := notmatches.Build()
		return GetMatchesRespOK(resp)
	}

	var userids []int64
	for _, val := range ent {
		if val.UserID == user.UserID {
			userids = append(userids, val.PartnerID)
			continue
		}
		userids = append(userids, val.UserID)
	}

	// Find all match user profile
	ents, _ := userHandler.FindByIDs(userids)
	var allmatches entities.MatchUserResponses
	for _, val := range ents {
		var tmp = entities.MatchUserResponse{}
		tmp.ApplyUser(val)

		allmatches = append(allmatches, tmp)
	}
	sEnt := allmatches.Build()

	return si.NewGetMatchesOK().WithPayload(sEnt)
}

//return 200
func GetMatchesRespOK(response []*models.MatchUserResponse) middleware.Responder {
	return si.NewGetMatchesOK().WithPayload(response)
}

// return 400 Bad Request
func GetMatchessRespBadReqestErr() middleware.Responder {
	return si.NewGetMatchesBadRequest().WithPayload(
		&si.GetMatchesBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

// return 401 Token Is Invalid
func GetMatchessRespUnauthErr() middleware.Responder {
	return si.NewGetMatchesUnauthorized().WithPayload(
		&si.GetMatchesUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

// return 500 Internal Server Error
func GetMatchesRespInternalErr() middleware.Responder {
	return si.NewGetMatchesInternalServerError().WithPayload(
		&si.GetMatchesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}
