package userlike

import (
	"fmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	nulr := repositories.NewUserLikeRepository()
	nutr := repositories.NewUserTokenRepository()
	nur := repositories.NewUserRepository()
	numr := repositories.NewUserMatchRepository()

	usrid, err := nutr.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "404",
				Message: "User Token Not Found",
			})
	}
	// find already matching user
	match, err := numr.FindAllByUserID(usrid.UserID)
	if err != nil {
		fmt.Println(err)
	}

	// find recive like except already matching user
	usrs, err := nulr.FindGotLikeWithLimitOffset(usrid.UserID, int(p.Limit), int(p.Offset), match)
	if err != nil {
		fmt.Println(err)
	}
	var userids []int64
	for i := 0; i < len(usrs); i++ {
		userids = append(userids, usrs[i].UserID)
	}
	//jimae
	ent, err := nur.FindByIDs(userids)
	if err != nil {
		fmt.Println(err)
	}
	var likeuserresp entities.LikeUserResponses
	for _, val := range ent {
		var tmp = entities.LikeUserResponse{}
		tmp.ApplyUser(val)
		likeuserresp = append(likeuserresp, tmp)
	}
	sEnt := likeuserresp.Build()

	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	return si.NewPostLikeOK()
}
