package usermatch

import (
	"fmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	usertokenHandler := repositories.NewUserTokenRepository()
	usermatchHandler := repositories.NewUserMatchRepository()
	userHandler := repositories.NewUserRepository()
	user, err := usertokenHandler.GetByToken(p.Token)
	if err != nil {
		return GetMatchessRespUnauthErr()
	}
	ent, err := usermatchHandler.FindByUserIDWithLimitOffset(user.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		return GetMatchesRespInternalErr()
	}
	var userids []int64
	for _, val := range ent {
		userids = append(userids, val.PartnerID)
	}
	fmt.Println(userids)
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
