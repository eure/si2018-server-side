package userlike

import (
	"fmt"

	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	nulr := repositories.NewUserLikeRepository()
	nutr := repositories.NewUserTokenRepository()
	//nur := repositories.NewUserRepository()
	numr := repositories.NewUserMatchRepository()

	usrid, err := nutr.GetByToken(p.Token)
	if err != nil {
		fmt.Println(err)
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "404",
				Message: "User Token Not Found",
			})
	}
	fmt.Println(usrid)
	match, err := numr.FindAllByUserID(usrid.UserID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(match)

	usrs, err := nulr.FindGotLikeWithLimitOffset(usrid.UserID, int(p.Limit), int(p.Offset), match)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(usrs)
	sample, _ := nulr.FindLikeAll(usrid.UserID)
	fmt.Println(sample)
	var userids []int64
	for i := 0; 0 < len(usrs); i++ {
		userids = append(userids, usrs[i].UserID)
	}
	ent, err := nulr.FindAllLikeUserResponse(sample)
	var ud []*models.LikeUserResponse
	for i := 0; i < len(ent); i++ {
		us := ent[i].Build()
		ud = append(ud, &us)
	}
	return si.NewGetLikesOK().WithPayload(ud)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	return si.NewPostLikeOK()
}
