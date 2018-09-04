package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/common"
)

// token -> userIDへ送信
func PostMessage(p si.PostMessageParams) middleware.Responder {
	ur := repositories.NewUserRepository()
	mr := repositories.NewUserMessageRepository()
	mchr := repositories.NewUserMatchRepository()
	// lr := repositories.NewUserLikeRepository()

	usr, _ := ur.GetByToken(p.Params.Token)

	matchIDs, _ := mchr.FindAllByUserID(usr.ID)

	if !common.Contains(matchIDs, p.UserID) {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code: "400",
				Message: "You are bocchi",
			})
	}
	msg := entities.NewUserMessage(usr.ID, p.UserID, p.Params.Message)
	mr.Create(msg)

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code: "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	mr := repositories.NewUserMessageRepository()
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
