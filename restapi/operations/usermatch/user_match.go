package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	r := repositories.NewUserMatchRepository()
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

	ent, err := r.FindByUserIDWithLimitOffset(loginUserToken.UserID,int(p.Limit),int(p.Offset))
	if err != nil {
		return outPutGetStatus(500)
	}
	if ent == nil {
		return outPutGetStatus(400)
	}
	
	// applied メソッドによって変換されたUsersがほしい。
	// とりあえずほしいから，格納先を用意してあげる。
	var appliedUsers entities.MatchUserResponses
	for _,m := range ent {
		var applied = entities.MatchUserResponse{}
		matchedUsers, err := u.GetByUserID(m.PartnerID)
		if err != nil {
			return outPutGetStatus(500)
		}

		applied.ApplyUser(*matchedUsers)
		appliedUsers = append (appliedUsers, applied)
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