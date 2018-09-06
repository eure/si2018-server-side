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
	userlikeHandler := repositories.NewUserLikeRepository()
	usertokenHandler := repositories.NewUserTokenRepository()
	userHandler := repositories.NewUserRepository()
	usermatchhandler := repositories.NewUserMatchRepository()
	token := p.Token
	limit := int(p.Limit)
	offset := int(p.Offset)
	usertoken, err := usertokenHandler.GetByToken(token)
	if err != nil {
		return GetLikesRespUnauthErr()
	}
	if usertoken == nil {
		return GetLiksRespBadReqestErr()
	}
	// find already matching user
	match, err := usermatchhandler.FindAllByUserID(usertoken.UserID)
	if err != nil {
		GetLikesRespInternalErr()
	}

	// find recive like except already matching user
	usrs, err := userlikeHandler.FindGotLikeWithLimitOffset(usertoken.UserID, limit, offset, match)
	if err != nil {
		GetLikesRespInternalErr()
	}
	var userids []int64
	for i := 0; i < len(usrs); i++ {
		userids = append(userids, usrs[i].UserID)
	}
	ent, err := userHandler.FindByIDs(userids)
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
	userlikeHandler := repositories.NewUserLikeRepository()
	usertokenHandler := repositories.NewUserTokenRepository()
	userHandler := repositories.NewUserRepository()
	// find myuser data
	userID := p.UserID
	postlikeParam := p.Params
	usertoken, err := usertokenHandler.GetByToken(postlikeParam.Token)
	if err != nil {
		return PosLikesRespInternalErr()
	}
	if usertoken == nil {
		return PostLiksRespBadReqestErr()
	}
	// validate if send same sex
	partner, err := userHandler.GetByUserID(userID)
	if err != nil {
		return PosLikesRespInternalErr()
	}
	user, err := userHandler.GetByUserID(usertoken.UserID)

	if err != nil {
		return PosLikesRespInternalErr()
	}
	if user.GetOppositeGender() == partner.GetOppositeGender() {
		return PostLiksRespBadReqestErr()
	}
	// duplicate like send
	duplicatelike, err := userlikeHandler.GetLikeBySenderIDReceiverID(usertoken.UserID, userID)
	if err != nil {
		PostLiksRespBadReqestErr()
	}
	if duplicatelike == nil {
		PostLiksRespBadReqestErr()
	}
	var userlike entities.UserLike
	BindUserLike(&userlike, userID, partner.ID)

	err = userlikeHandler.Create(userlike)
	if err != nil {
		PosLikesRespInternalErr()
	}
	fmt.Println(err)
	return PostLikeOK()
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

func PostLikeOK() middleware.Responder {
	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:    "200",
			Message: "OK",
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
