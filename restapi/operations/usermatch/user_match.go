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
	loginUserToken , _ := t.GetByToken(token)

	ent, _ := r.FindByUserIDWithLimitOffset(loginUserToken.UserID,int(p.Limit),int(p.Offset))

	// applied メソッドによって変換されたUsersがほしい。
	// とりあえずほしいから，格納先を用意してあげる。
	var appliedusers entities.MatchUserResponses
	for _,m := range ent {
		var applied = entities.MatchUserResponse{}
		matcheduser , _ := u.GetByUserID(m.PartnerID)
		applied.ApplyUser(*matcheduser)
		appliedusers = append (appliedusers, applied)
	}

	sEtc := appliedusers.Build()
	return si.NewGetMatchesOK().WithPayload(sEtc)
}
