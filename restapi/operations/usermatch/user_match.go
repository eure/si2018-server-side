package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

//- マッチング一覧API
//- GET {hostname}/api/1.0/matches
//- ページネーションを実装してください
//- TokenのValidation処理を実装してください
//- ※お互いいいね！を送るとmatchingになります
func GetMatches(p si.GetMatchesParams) middleware.Responder {
	repUserMatch := repositories.NewUserMatchRepository()
	repUserToken := repositories.NewUserTokenRepository() // トークンからユーザのIDを取得するため
	repUser := repositories.NewUserRepository() // IDからUserを取得するため

	// tokenのバリデーション
	err := repUserToken.ValidateToken(p.Token)
	if err != nil {
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code: "401",
				Message: "Token Is Invalid",
			})
	}

	// bad Request
	if (p.Limit < 1) || (p.Offset < 0) {
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// トークンからidの取得
	userToken, _ := repUserToken.GetByToken(p.Token)

	// マッチング済みの相手一覧を取得する.
	userMatches, err := repUserMatch.FindByUserIDWithLimitOffset(userToken.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userMatches == nil {
		return si.NewGetMatchesOK().WithPayload(nil) // UserMatchがいない場合, nilで返す
	}

	// マッチングしているユーザのIDを配列にいれる
	var matchedUserIDs []int64
	var matchedUserID int64
	for _, matchUser:= range userMatches {
		// UserMatchのUserID, PartnerIDのどちらがマッチングしている相手のIDか調べる。
		if matchUser.UserID == userToken.UserID {
			matchedUserID = matchUser.PartnerID
		} else {
			matchedUserID = matchUser.UserID
		}
		matchedUserIDs = append(matchedUserIDs, matchedUserID)
	}

	// IDの配列からユーザーを取得
	matchedUsers, err := repUser.FindByIDs(matchedUserIDs)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// return用のmodelを作るためのMatchUserResponsesのエンティティ
	var matchUserReses entities.MatchUserResponses

	for i, matchedUser := range matchedUsers {
		// MatchUserResponsesに入れていくためのMatchUserResponseのエンティティの実態を宣言
		var matchUserRese = entities.MatchUserResponse{}
		// User型からMatchUserResponse型に変換
		matchUserRese.ApplyUser(matchedUser)
		matchUserRese.MatchedAt = userMatches[i].CreatedAt
		matchUserReses = append(matchUserReses, matchUserRese)
	}

	matchUserResesModel := matchUserReses.Build()

	return si.NewGetMatchesOK().WithPayload(matchUserResesModel)
}
