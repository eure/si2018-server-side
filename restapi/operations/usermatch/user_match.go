package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	rt := repositories.NewUserTokenRepository()
	r := repositories.NewUserRepository()
	rm := repositories.NewUserMatchRepository()
	t, err := rt.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByToken failed: " + err.Error(),
			})
	}
	if t == nil {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗): GetByToken failed",
			})
	}
	matched, err := rm.FindByUserIDWithLimitOffset(t.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: FindAllByUserID failed: " + err.Error(),
			})
	}
	ids := make([]int64, 0)
	for _, m := range matched {
		ids = append(ids, m.GetPartnerID(t.UserID))
	}
	users, err := r.FindByIDs(ids)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: FindByIDs failed: " + err.Error(),
			})
	}
	sEnt := make([]*models.MatchUserResponse, 0)
	for i, m := range matched {
		response := entities.MatchUserResponse{MatchedAt: m.CreatedAt}
		response.ApplyUser(users[i])
		swaggerMatch := response.Build()
		sEnt = append(sEnt, &swaggerMatch)
	}
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
