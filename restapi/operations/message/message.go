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
	tr := repositories.NewUserTokenRepository()
	ur := repositories.NewUserRepository()
	mr := repositories.NewUserMessageRepository()
	mchr := repositories.NewUserMatchRepository()

	token, err := tr.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error (in get token)",
			})
	}

	//
	if token == nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code: "401",
				Message: "Unauthorized",
			})
	}

	usr, err := ur.GetByUserID(token.UserID)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error(in get user)",
			})
	}

	if usr == nil {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	matchIDs, err := mchr.FindAllByUserID(usr.ID)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	if !common.Contains(matchIDs, p.UserID) {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code: "400",
				Message: "Bad Request (not match)",
			})
	}

	msg := entities.NewUserMessage(usr.ID, p.UserID, p.Params.Message)
	err = mr.Create(msg)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code: "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	tr := repositories.NewUserTokenRepository()
	mr := repositories.NewUserMessageRepository()
	matchr := repositories.NewUserMatchRepository()

	token, err := tr.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code: "500",
				Message: "ISE (in token)",
			})
	}

	if token == nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code: "401",
				Message: "Unauthorized",
			})
	}

	partnerIDs, err := matchr.FindAllByUserID(p.UserID)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code: "500",
				Message: "ISE (in partner ids)",
			})
	}
	if partnerIDs == nil {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code: "400",
				Message: "Bad Request (parnerIDs is nil)",
			})
	}


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
