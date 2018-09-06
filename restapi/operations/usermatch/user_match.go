package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	m := repositories.NewUserMatchRepository()
	u := repositories.NewUserRepository()

	// ユーザーID取得用
	token, err := t.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if token == nil {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code:    "401",
				Message: "Your Token Is Invalid",
			})
	}

	// マッチ済みのユーザーを取得
	match, err := m.FindByUserIDWithLimitOffset(token.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if match == nil {
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// ApplyUserの形にしてBuildしないといけないので
	// ApplyUserの型に変換する
	var mr entities.MatchUserResponses
	for _, matchUser := range match {
		// 構造体の初期化(単体)
		r := entities.MatchUserResponse{}

		// Userの情報を取得
		user, err := u.GetByUserID(matchUser.PartnerID)
		if err != nil {
			return si.NewGetMatchesInternalServerError().WithPayload(
				&si.GetMatchesInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
		if user == nil {
			return si.NewGetMatchesBadRequest().WithPayload(
				&si.GetMatchesBadRequestBody{
					Code:    "400",
					Message: "Bad Request",
				})
		}

		r.ApplyUser(*user)
		mr = append(mr, r)
	}

	sEnt := mr.Build()
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
