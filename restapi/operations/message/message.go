package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	return si.NewPostMessageOK()
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	mr := repositories.NewUserMessageRepository()
	// ur := repositories.NewUserRepository()
	matchr := repositories.NewUserMatchRepository()
	partnerIDs, _ := matchr.FindAllByUserID(p.UserID)
	var msgses entities.UserMessages
	for _, pID := range partnerIDs {
		messages, err := mr.GetMessages(p.UserID, pID, int(*p.Limit), p.Latest, p.Oldest)
		if err != nil {
			return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
		}
		for _, msg := range messages{
			msgses = append(msgses, msg)
		}
	}

	sMsgs := msgses.Build()

	return si.NewGetMessagesOK().WithPayload(sMsgs)
}
