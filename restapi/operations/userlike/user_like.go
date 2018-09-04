package userlike

import (
	"fmt"
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
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
	nulr := repositories.NewUserLikeRepository()
	nutr := repositories.NewUserTokenRepository()
	//numr := repositories.NewUserMatchRepository()
	nur := repositories.NewUserRepository()
	// find myuser data
	userid, _ := nutr.GetByToken(p.Params.Token)
	// validate if send same sex
	partner, err := nur.GetByUserID(p.UserID)
	if err != nil {
		fmt.Println(err)
	}
	user, err := nur.GetByUserID(userid.UserID)
	if err != nil {
		fmt.Println(err)
	}
	if user.GetOppositeGender() == partner.GetOppositeGender() {
		fmt.Println("err")
	}
	// duplicate like send
	duplicatelike, err := nulr.GetLikeBySenderIDReceiverID(userid.UserID, p.UserID)
	if err != nil {
		fmt.Println(err)
	}
	if duplicatelike == nil {
		fmt.Println(duplicatelike)
	}
	var userlike entities.UserLike
	BindUserLike(&userlike, userid.UserID, partner.ID)

	err = nulr.Create(userlike)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("OK")
	return si.NewPostLikeOK()
}
func BindUserLike(like *entities.UserLike, userid int64, partnerid int64) {
	like.UserID = userid
	like.PartnerID = partnerid
	like.CreatedAt = strfmt.DateTime(time.Now())
	like.UpdatedAt = strfmt.DateTime(time.Now())
}
