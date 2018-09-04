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

	userToken, _ := repUserToken.GetByToken(p.Token)

	// マッチング済みの相手一覧を取得する.
	UserMatches, err := repUserMatch.FindByUserIDWithLimitOffset(userToken.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if UserMatches == nil {
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code:    "400",
				Message: "Nobody matched.",
			})
	}

	// マッチングしているユーザのIDを配列にいれる
	var matchedUserIDs []int64
	var matchedUserID int64
	for _, matchUser:= range UserMatches {
		// UserMatchのUserID, PartnerIDのどちらがマッチングしている相手のIDか調べる。
		if matchUser.UserID == userToken.UserID {
			matchedUserID = matchUser.PartnerID
		} else {
			matchedUserID = matchUser.UserID
		}
		matchedUserIDs = append(matchedUserIDs, matchedUserID)
	}

	// IDの配列からユーザーを取得
	matchedUsers, _ := repUser.FindByIDs(matchedUserIDs)

	// return用のmodelを作るためのMatchUserResponsesのエンティティ
	var matchUserReses entities.MatchUserResponses

	for _, matchedUser := range matchedUsers {
		var matchUserRese = entities.MatchUserResponse{}
		matchUserRese.ApplyUser(matchedUser)
		matchUserReses = append(matchUserReses, matchUserRese)
	}

	matchUserResesModel := matchUserReses.Build()

	return si.NewGetMatchesOK().WithPayload(matchUserResesModel)
}
