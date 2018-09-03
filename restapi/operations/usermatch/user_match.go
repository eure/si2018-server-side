package usermatch

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"

	"github.com/eure/si2018-server-side/entities"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	// TODO: 400エラー
	// TODO: 401エラー
	// tokenからUserIDを取得
	tokenR         := repositories.NewUserTokenRepository()
	tokenEnt, err  := tokenR.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	// マッチしているユーザの取得(IDしか取れない？)※ UserID, PartnerID, CreatedAt, UpdatedAt
	// なのでこの後でPartnerIDを使用してマッチングしているユーザの情報を取得する必要があると考えた
	matchR              := repositories.NewUserMatchRepository()
	matchUsers, err := matchR.FindByUserIDWithLimitOffset(tokenEnt.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	// マッチしているuserIdsからユーザ情報の取得
	var matchUserIds []int64
	for _, u := range matchUsers {
		matchUserIds = append(matchUserIds, u.PartnerID)
	}
	userR := repositories.NewUserRepository()
	responseModels, err := userR.FindByIDs(matchUserIds)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	var array entities.MatchUserResponses
	for _, u := range responseModels {
		var tmp = entities.MatchUserResponse{}
		tmp.ApplyUser(u)
		array = append(array, tmp)
	}

	return si.NewGetMatchesOK().WithPayload(array.Build())
}
