package usermatch

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {

	// レポジトリを初期化する
	tokenR := repositories.NewUserTokenRepository()
	userMatchR := repositories.NewUserMatchRepository()
	userR := repositories.NewUserRepository()

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Token)

	// 401エラー
	if tokenEnt == nil {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code: "401",
				Message: "Your token is invalid.",
			})
	}

	// 500エラー
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	// マッチのデータを取得する
	userMatchEnts, err := userMatchR.FindByUserIDWithLimitOffset(tokenEnt.UserID, int(p.Limit), int(p.Offset))

	// 500エラー
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	// マッチしているユーザーIDのスライスを生成する
	var matchUserIds []int64

	for _, userMatchEnt := range userMatchEnts {

		if tokenEnt.UserID == userMatchEnt.UserID {

			matchUserIds = append(matchUserIds, userMatchEnt.PartnerID)

		}

		if tokenEnt.UserID == userMatchEnt.PartnerID {

			matchUserIds = append(matchUserIds, userMatchEnt.UserID)

		}

	}

	// マッチしているユーザー情報を取得する
	userEnts, err := userR.FindByIDs(matchUserIds)

	// 500エラー
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}


	// レスポンス用の構造体に必要な情報をマップする
	var getMatchesResponsesEnt entities.MatchUserResponses

	for _, userMatchEnt := range userMatchEnts {

		var getMatchesResponseEnt entities.MatchUserResponse

		getMatchesResponseEnt.MatchedAt = userMatchEnt.CreatedAt

		for _, userEnt := range userEnts {

			if userMatchEnt.PartnerID == userEnt.ID {
				getMatchesResponseEnt.ApplyUser(userEnt)
			}

			if userMatchEnt.UserID == userEnt.ID {
				getMatchesResponseEnt.ApplyUser(userEnt)
			}

		}

		getMatchesResponsesEnt = append(getMatchesResponsesEnt, getMatchesResponseEnt)

	}

	// モデルに変換する
	getMatchesresponses := getMatchesResponsesEnt.Build()

	// 結果を返す
	return si.NewGetMatchesOK().WithPayload(getMatchesresponses)
}
