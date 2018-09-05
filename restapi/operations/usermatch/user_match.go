package usermatch

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	tr := repositories.NewUserTokenRepository()
	mr := repositories.NewUserMatchRepository()
	ur := repositories.NewUserRepository()

	token, err := tr.GetByToken(p.Token)
	if err != nil {
		si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code: "500",
				Message: "ISE (in get token)",
			})
	}

	if token == nil {
		si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code: "401",
				Message: "Unauthorized",
			})
	}

	user, err := ur.GetByUserID(token.UserID)
	if err != nil {
		si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code: "500",
				Message: "ISE (in get token)",
			})
	}
	if user == nil {
		si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	var ents entities.UserMatches
	ents, err = mr.FindByUserIDWithLimitOffset(user.ID, int(p.Limit), int(p.Offset))
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	var responses entities.MatchUserResponses

	for _, matched := range ents {
		res := entities.MatchUserResponse{}
		user, err := ur.GetByUserID(matched.UserID)
		if err != nil {
			return si.NewGetMatchesInternalServerError().WithPayload(
				&si.GetMatchesInternalServerErrorBody{
					Code: "500",
					Message: "Internal Server Error",
				})
		}

		if 
		res.ApplyUser(*user)
		responses = append(responses, res)
	}

	sResponses := responses.Build()

	return si.NewGetMatchesOK().WithPayload(sResponses)
}
