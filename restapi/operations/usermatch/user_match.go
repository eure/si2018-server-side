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
				Message: "Token Is Invalid",
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

	// 明示的に型宣言
	var mr entities.MatchUserResponses
	var userIDs []int64

	for _, matchUser := range match {
		userIDs = append(userIDs, matchUser.UserID)
	}
	// マッチしているユーザーの情報をすべて取得する
	users, err := u.FindByIDs(userIDs)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// ApplyUserの形にしてBuildしないといけないので
	// ApplyUserの型に変換する
	for _, user := range users {
		// 構造体の初期化
		r := entities.MatchUserResponse{}
		r.ApplyUser(user)
		mr = append(mr, r)
	}

	sEnt := mr.Build()
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
