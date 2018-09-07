package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	m := repositories.NewUserMatchRepository()
	t := repositories.NewUserTokenRepository()
	u := repositories.NewUserRepository()
	
	// tokenから UserToken entitiesの取得 (Validation)
	token := p.Token
	loginUserToken , err := t.GetByToken(token)
	if err != nil {
		return outPutGetStatus(500)
	}
	if loginUserToken == nil {
		return outPutGetStatus(401)
	}
	
	// limit が20かどうか検出
	if p.Limit != int64(20) {
		return outPutGetStatus(400)
	}
	// offset が0以上かどうか検出
	if p.Offset >= int64(0) {
		return outPutGetStatus(400)
	}

	ent, err := m.FindByUserIDWithLimitOffset(loginUserToken.UserID,int(p.Limit),int(p.Offset))
	if err != nil {
		return outPutGetStatus(500)
	}
	if ent == nil {
		return outPutGetStatus(400)      // <<<<<<<<<　ここは，マッチしてる人がいなかったら空のJSONでいいの？
	}
	
	// applied メソッドによって変換されたUsersがほしい。
	// とりあえずほしいから，格納先を用意してあげる。
	var appliedUsers entities.MatchUserResponses
	for _,m := range ent {
		var applied = entities.MatchUserResponse{}
		if m.PartnerID == loginUserToken.UserID {
			matchedUsers, err := u.GetByUserID(m.UserID)
			if err != nil {
				return outPutGetStatus(500)
			}
			applied.ApplyUser(*matchedUsers)
			appliedUsers = append (appliedUsers, applied)
		} else {
			matchedUsers, err := u.GetByUserID(m.PartnerID)
			if err != nil {
				return outPutGetStatus(500)
			}
			applied.ApplyUser(*matchedUsers)
			appliedUsers = append (appliedUsers, applied)
		}
	}

	sEtc := appliedUsers.Build()
	return si.NewGetMatchesOK().WithPayload(sEtc)
}


func outPutGetStatus(num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	case 401:
		return si.NewGetMatchesUnauthorized().WithPayload(
			&si.GetMatchesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗)",
			})
	case 400:
		return si.NewGetMatchesBadRequest().WithPayload(
			&si.GetMatchesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	return nil
}