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

func PostMessage(p si.PostMessageParams) middleware.Responder { /* TODO 連打 */
	message := p.Params.Message
	token := p.Params.Token
	pid := p.UserID

	// Validations
	err1 := util.ValidateLimit(limit)
	err2 := util.ValidateOffset(offset)
	if (err1 != nil) || (err2 != nil) {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	err := util.ValidateToken(token)
	if err != nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	uid, _ := util.GetIDByToken(token)

	// Already matching?
	rm := repositories.NewUserMatchRepository()

	mat, err := rm.Get(uid, pid)
	if err != nil {
		fmt.Print("Get err: ")
		fmt.Println(err)
		return si.NewPostMessagesInternalServerError().WithPayload(
			&si.PostMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if mat == nil {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// Prepare new UserMessage
	ent := entities.UserMessage{}
	ent.UserID = uid
	ent.PartnerID = pid
	ent.Message = message
	ent.CreatedAt = strfmt.DateTime(time.Now())
	ent.UpdatedAt = strfmt.DateTime(time.Now())

	// Add
	rs := repositories.NewUserMessageRepository()
	err = rs.Create(ent)
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
	/* TODO check [lat|old]est? */

	// Validations
	err1 := util.ValidateLimit(limit)
	err2 := util.ValidateOffset(offset)
	if (err1 != nil) || (err2 != nil) {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	err := util.ValidateToken(token)
	if err != nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	uid, _ := util.GetIDByToken(token)
	
	// Already matching?
	rm := repositories.NewUserMatchRepository()

	mat, err := rm.Get(uid, pid)
	if err != nil {
		fmt.Print("Get err: ")
		fmt.Println(err)
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if mat == nil {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// Get messages in DESC order
	rs := repositories.NewUserMessageRepository()

	messages, err := rs.GetMessages(uid, pid, int(limit), latest, oldest)
	if err != nil {
		fmt.Print("Get messages err: ")
		fmt.Println(err)
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// Prepare response
	var reses entities.UserMessages
	reses = messages

	return si.NewGetMessagesOK().WithPayload(reses.Build())
}
