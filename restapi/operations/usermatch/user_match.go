package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	r := repositories.NewUserMatchRepository()

	ruserR := repositories.NewUserRepository()
	user , _ := ruserR.GetByToken(p.Token)
	userid := user.ID

	ent, _ := r.FindByUserIDWithLimitOffset(userid,int(p.Limit),int(p.Offset))

	// applied メソッドによって変換されたUser's'がほしい。
	// とりあえずほしいから，格納先を用意してあげる。
	var appliedusers entities.MatchUserResponses
	for _,m := range ent {
		var applied = entities.MatchUserResponse{}
		matcheduser , _ := ruserR.GetByUserID(m.UserID)
		applied.ApplyUser(*matcheduser)
		appliedusers = append (appliedusers, applied)
	}

	sEtc := appliedusers.Build()
	return si.NewGetMatchesOK().WithPayload(sEtc)
}
