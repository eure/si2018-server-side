package usermatch

import (
	"fmt"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	nutr := repositories.NewUserTokenRepository()
	numr := repositories.NewUserMatchRepository()
	user, err := nutr.GetByToken(p.Token)
	ent, err := numr.FindByUserIDWithLimitOffset(user.UserID, int(p.Limit), int(p.Offset))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ent)

	return si.NewGetMatchesOK()
}
