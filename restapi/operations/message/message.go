package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/restapi/operations/util"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/runtime/middleware"
	"fmt"
	"time"
	"github.com/go-openapi/strfmt"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	message := p.Params.Message
	token := p.Params.Token
	pid := p.UserID
	/* TODO bad request */

	err := util.ValidateToken(token)
	if err != nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	uid, _ := util.GetIDByToken(token)

	rl := repositories.NewUserLikeRepository()
	like, err := rl.GetLikeBySenderIDReceiverID(uid, pid)
	if err != nil {
		fmt.Print("Get likes err: ")
		fmt.Println(err)
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if like == nil {
		fmt.Println("Not matching yet")
		return si.NewPostMessageBadRequest().WithPayload( /* TODO 403? */
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request (Not matching yet)",
			})
	}

	ent := entities.UserMessage{}
	ent.UserID = uid
	ent.PartnerID = pid
	ent.Message = message
	ent.CreatedAt = strfmt.DateTime(time.Now())
	ent.UpdatedAt = strfmt.DateTime(time.Now())

	rm := repositories.NewUserMessageRepository()
	err = rm.Create(ent)
	if err != nil {
		fmt.Print("Create message err: ")
		fmt.Println(err)
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	
	return si.NewPostMessageOK()
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	token := p.Token
	pid := p.UserID
	latest := p.Latest
	oldest := p.Oldest
	limit := *p.Limit
	/* TODO bad request */

	err := util.ValidateToken(token)
	if err != nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	uid, _ := util.GetIDByToken(token)
	
	/* TODO validate matching state? */

	r := repositories.NewUserMessageRepository()

	messages, err := r.GetMessages(uid, pid, int(limit), latest, oldest) /* TODO order */
	if err != nil {
		fmt.Print("Get messages err: ")
		fmt.Println(err)
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	var reses entities.UserMessages
	reses = messages

	return si.NewGetMessagesOK().WithPayload(reses.Build())
}
