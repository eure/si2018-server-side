package message

import (
	"fmt"
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
	usertoken, err := nutr.GetByToken(p.Params.Token)
	if err != nil {
		fmt.Println(err)
	}

	// validate already matching
	existmatch, err := numar.Get(usertoken.UserID, p.UserID)
	if err != nil {
		fmt.Println(err)
	}
	if existmatch == nil {
		fmt.Println(err)
	}

	user := entities.UserMessage{
		UserID:    usertoken.UserID,
		PartnerID: p.UserID,
		Message:   p.Params.Message,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}

	numr.Create(user)

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	numr := repositories.NewUserMessageRepository()
	nutr := repositories.NewUserTokenRepository()
	nur := repositories.NewUserRepository()
	usertoken, err := nutr.GetByToken(p.Token)
	if err != nil {
		fmt.Println(err)
	}
	userid, err := nur.GetByUserID(p.UserID)
	if err != nil {
		fmt.Println(err)
	}
	msg, _ := numr.GetMessages(usertoken.UserID, userid.ID, int(*p.Limit), p.Latest, p.Oldest)
	var respmsg entities.UserMessages
	for _, msg := range msg {
		respmsg = append(respmsg, msg)

	}
	sEnt := respmsg.Build()
	return si.NewGetMessagesOK().WithPayload(sEnt)
}
