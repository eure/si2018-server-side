package usermatch

import (
	"fmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	nutr := repositories.NewUserTokenRepository()
	numr := repositories.NewUserMatchRepository()
	nur := repositories.NewUserRepository()
	user, err := nutr.GetByToken(p.Token)
	ent, err := numr.FindByUserIDWithLimitOffset(user.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		fmt.Println(err)
	}
	var userids []int64
	for _, val := range ent {
		userids = append(userids, val.PartnerID)
	}
	fmt.Println(userids)
	ents, _ := nur.FindByIDs(userids)
	var allmatches entities.MatchUserResponses
	for _, val := range ents {
		var tmp = entities.MatchUserResponse{}
		tmp.ApplyUser(val)
		allmatches = append(allmatches, tmp)
	}
	sEnt := allmatches.Build()

	return si.NewGetMatchesOK().WithPayload(sEnt)
}
