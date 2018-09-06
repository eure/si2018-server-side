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

	// paramsの変数を定義
	paramsToken := p.Token
	paramsLimit := p.Limit
	paramsOffset := p.Offset

	// ユーザーID取得用
	token, err := t.GetByToken(paramsToken)
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

	// imitが20になっているかをvalidation
	if paramsLimit != int64(20) {
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// offset 負の数になったときにBadRequestを返す
	if paramsOffset < 0 {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// マッチ済みのユーザーを取得
	match, err := m.FindByUserIDWithLimitOffset(token.UserID, int(paramsLimit), int(paramsOffset))
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 明示的に型宣言
	var mr entities.MatchUserResponses
	var userIDs []int64

	// マッチ済みのユーザーのIDをスライスに入れる
	for _, matchUser := range match {
		if matchUser.UserID != token.UserID {
			userIDs = append(userIDs, matchUser.UserID)
		} else {
			userIDs = append(userIDs, matchUser.PartnerID)
		}
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
	for i := len(users) - 1; i >= 0; i-- {
		// 構造体の初期化
		r := entities.MatchUserResponse{}
		r.ApplyUser(users[i])
		mr = append(mr, r)
	}

	sEnt := mr.Build()
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
