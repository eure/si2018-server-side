package usermatch

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"

	"github.com/eure/si2018-server-side/entities"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	// LIMIT OFFSET Check -> 400
	if p.Limit <= 0 || p.Offset < 0 {
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code   : "400",
				Message: "Bad Request",
			})
	}
	// Tokenのユーザが存在しない -> 401
	tokenR        := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	if tokenEnt == nil{
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code   : "401",
				Message: "Token Is Invalid",
			})
	}

	// tokenからUserIDを取得
	// マッチしているユーザの取得(IDしか取れない？)※ UserID, PartnerID, CreatedAt, UpdatedAt
	// なのでこの後でPartnerIDを使用してマッチングしているユーザの情報を取得する必要があると考えた
	matchR            := repositories.NewUserMatchRepository()
	matchEntList, err := matchR.FindByUserIDWithLimitOffset(tokenEnt.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	// マッチしているユーザIDを作成
	var matchIds []int64
	for _, u := range matchEntList {
		matchIds = append(matchIds, u.PartnerID)
	}
	// マッチしているユーザIDからPartnerユーザ情報の取得
	userR            := repositories.NewUserRepository()
	userEntList, err := userR.FindByIDs(matchIds)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	var array entities.MatchUserResponses
	for _, u := range userEntList {
		var tmp entities.MatchUserResponse
		tmp.ApplyUser(u)
		array = append(array, tmp)
	}
	responseData := array.Build()
	return si.NewGetMatchesOK().WithPayload(responseData)
}
