package message

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	return si.NewPostMessageOK()
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	// XXX このトークン認証が何回もありすぎてクソイケてない
	userTokenEnt, err := repositories.NewUserTokenRepository().GetByToken(p.Token)

	if err != nil {
		return getMessagesInternalServerErrorResponse()
	}

	if userTokenEnt == nil {
		return getMessageUnauthorizedResponse()
	}

	// int64になっているのでcastする必要がある
	limit := int(*p.Limit)
	userID := userTokenEnt.UserID
	partnerID := p.UserID
	latest := p.Latest
	oldest := p.Oldest

	userMessageRepository := repositories.NewUserMessageRepository()

	var userMessagesEnt entities.UserMessages
	userMessagesEnt, err = userMessageRepository.GetMessages(userID, partnerID, limit, latest, oldest)

	userMessages := userMessagesEnt.Build()

	return si.NewGetMessagesOK().WithPayload(userMessages)
}

func getMessagesInternalServerErrorResponse() middleware.Responder {
	return si.NewGetMessagesInternalServerError().WithPayload(
		&si.GetMessagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getMessageUnauthorizedResponse() middleware.Responder {
	return si.NewGetMessagesUnauthorized().WithPayload(
		&si.GetMessagesUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func postMessageOKResponse(message string) middleware.Responder {
	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: message,
		})
}

func postMessageUnauthorizedResponse() middleware.Responder {
	return si.NewPostMessageUnauthorized().WithPayload(
		&si.PostMessageUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func postMessageInternalServerErrorResponse() middleware.Responder {
	return si.NewPostMessageInternalServerError().WithPayload(
		&si.PostMessageInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}
