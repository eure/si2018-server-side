package message

import (
	"fmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {

	return si.NewPostMessageOK()
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
