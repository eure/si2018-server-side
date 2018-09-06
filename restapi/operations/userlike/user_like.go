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
		return GetLikesRespUnauthErr()
	}
	// find already matching user
	match, err := numr.FindAllByUserID(usrid.UserID)
	if err != nil {
		GetLikesRespInternalErr()
	}

	// find recive like except already matching user
	usrs, err := nulr.FindGotLikeWithLimitOffset(usrid.UserID, int(p.Limit), int(p.Offset), match)
	if err != nil {
		GetLikesRespInternalErr()
	}
	var userids []int64
	for i := 0; i < len(usrs); i++ {
		userids = append(userids, usrs[i].UserID)
	}
	//jimae
	ent, err := nur.FindByIDs(userids)
	if err != nil {
		GetLikesRespInternalErr()
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
		return PosLikesRespInternalErr()
	}
	user, err := nur.GetByUserID(userid.UserID)
	if err != nil {
		return PosLikesRespInternalErr()
	}
	if user.GetOppositeGender() == partner.GetOppositeGender() {
		return PostLiksRespBadReqestErr()
	}
	// duplicate like send
	duplicatelike, err := nulr.GetLikeBySenderIDReceiverID(userid.UserID, p.UserID)
	if err != nil {
		PostLiksRespBadReqestErr()
	}
	if duplicatelike == nil {
		PostLiksRespBadReqestErr()
	}
	var userlike entities.UserLike
	BindUserLike(&userlike, userid.UserID, partner.ID)

	err = nulr.Create(userlike)
	if err != nil {
		PosLikesRespInternalErr()
	}
	fmt.Println(err)
	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:    "200",
			Message: "OK",
		})
}
func BindUserLike(like *entities.UserLike, userid int64, partnerid int64) {
	like.UserID = userid
	like.PartnerID = partnerid
	like.CreatedAt = strfmt.DateTime(time.Now())
	like.UpdatedAt = strfmt.DateTime(time.Now())
}

func GetLiksRespBadReqestErr() middleware.Responder {
	return si.NewGetLikesBadRequest().WithPayload(
		&si.GetLikesBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

func GetLikesRespUnauthErr() middleware.Responder {
	return si.NewGetLikesUnauthorized().WithPayload(
		&si.GetLikesUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

func GetLikesRespInternalErr() middleware.Responder {
	return si.NewGetLikesInternalServerError().WithPayload(
		&si.GetLikesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func PostLiksRespBadReqestErr() middleware.Responder {
	return si.NewPostLikeBadRequest().WithPayload(
		&si.PostLikeBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

func PostLikesRespUnauthErr() middleware.Responder {
	return si.NewPostLikeUnauthorized().WithPayload(
		&si.PostLikeUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

func PosLikesRespInternalErr() middleware.Responder {
	return si.NewPostLikeInternalServerError().WithPayload(
		&si.PostLikeInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}
