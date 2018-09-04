package usermatch

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"

	"github.com/eure/si2018-server-side/entities"
	"strings"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	// Tokenの形式がおかしい -> 401
	if !(strings.HasPrefix(p.Token, "USERTOKEN"))  {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code   : "401",
				Message: "Token Is Invalid",
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

	// マッチしているuserIdsからPartnerユーザ情報の取得
	var matchIds []int64
	for _, u := range matchEntList {
		matchIds = append(matchIds, u.PartnerID)
	}
	userR            := repositories.NewUserRepository()
	userEntList, err := userR.FindByIDs(matchIds)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	// TODO: 修正の余地あり?
	var array entities.MatchUserResponses
	for _, u := range userEntList {
		var tmp entities.MatchUserResponse
		tmp.ApplyUser(u)
		array = append(array, tmp)
	}

	responseData := array.Build()
	return si.NewGetMatchesOK().WithPayload(responseData)
}
