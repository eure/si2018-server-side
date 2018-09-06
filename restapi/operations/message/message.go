package message

import (
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	numr := repositories.NewUserMessageRepository()
	nutr := repositories.NewUserTokenRepository()
	//nur := repositories.NewUserRepository()
	numar := repositories.NewUserMatchRepository()
	// get my user id
	postParams := p.Params
	userID := p.UserID

	usertoken, err := nutr.GetByToken(postParams.Token)
	if err != nil {
		return PostMsgRespUnauthErr()
	}

	// Is There a collect user token?
	if usertoken == nil {
		return PostMsgBadReqestErr()
	}

	// Is there already matching?
	existmatch, err := numar.Get(usertoken.UserID, userID)
	if err != nil {
		return PostMsgRespInternalErr()
	}

	if existmatch == nil {
		return PostMsgBadReqestErr()
	}
	// validate already matching opposite
	existmatchopposite, err := numar.Get(usertoken.UserID, userID)
	if err != nil {
		return PostMsgRespInternalErr()
	}
	if existmatchopposite == nil {
		return PostMsgBadReqestErr()
	}

	user := entities.UserMessage{
		UserID:    usertoken.UserID,
		PartnerID: p.UserID,
		Message:   p.Params.Message,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}

	err = numr.Create(user)
	if err != nil {
		return PostMsgRespInternalErr()
	}

	return PostMsgOK()
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	numr := repositories.NewUserMessageRepository()
	nutr := repositories.NewUserTokenRepository()
	nur := repositories.NewUserRepository()
	usertoken, err := nutr.GetByToken(p.Token)
	if err != nil {
		return GetMsgRespUnauthErr()
	}
	userid, err := nur.GetByUserID(p.UserID)
	if err != nil {
		return GetMsgRespInternalErr()
	}
	msg, _ := numr.GetMessages(usertoken.UserID, userid.ID, int(*p.Limit), p.Latest, p.Oldest)
	var respmsg entities.UserMessages
	for _, msg := range msg {
		respmsg = append(respmsg, msg)

	}
	sEnt := respmsg.Build()
	return si.NewGetMessagesOK().WithPayload(sEnt)
}

func GetMsgBadReqestErr() middleware.Responder {
	return si.NewGetMessagesBadRequest().WithPayload(
		&si.GetMessagesBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

func GetMsgRespUnauthErr() middleware.Responder {
	return si.NewGetMessagesUnauthorized().WithPayload(
		&si.GetMessagesUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

func GetMsgRespInternalErr() middleware.Responder {
	return si.NewGetMessagesInternalServerError().WithPayload(
		&si.GetMessagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func PostMsgOK() middleware.Responder {
	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

func PostMsgBadReqestErr() middleware.Responder {
	return si.NewPostMessageBadRequest().WithPayload(
		&si.PostMessageBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

func PostMsgRespUnauthErr() middleware.Responder {
	return si.NewPostMessageUnauthorized().WithPayload(
		&si.PostMessageUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

func PostMsgRespInternalErr() middleware.Responder {
	return si.NewPostMessageInternalServerError().WithPayload(
		&si.PostMessageInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}
