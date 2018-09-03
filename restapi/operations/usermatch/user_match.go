package usermatch

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	mr := repositories.NewUserMatchRepository()
	ur := repositories.NewUserRepository()
	user, _ := ur.GetByToken(p.Token)
	var ents entities.UserMatches
	ents, err := mr.FindByUserIDWithLimitOffset(user.ID, int(p.Limit), int(p.Offset))
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
		user, _ := ur.GetByUserID(matched.UserID)
		res.ApplyUser(*user)
		responses = append(responses, res)
	}

	sResponses := responses.Build()

	return si.NewGetMatchesOK().WithPayload(sResponses)
}
