package message

import (
	"fmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	userTokenEnt, err := repositories.NewUserTokenRepository().GetByToken(p.Params.Token)

	if err != nil {
		return GetMessagesInternalServerErrorResponse()
	}

	if userTokenEnt == nil {
		return GetMessageUnauthorizedResponse()
	}

	userID := userTokenEnt.UserID
	partnerID := p.UserID
	message := p.Params.Message

	userMessageEnt := entities.UserMessage{
		UserID:    userID,
		PartnerID: partnerID,
		Message:   message,
	}

	messageRepository := repositories.NewUserMessageRepository()

	errs := messageRepository.Validate(userMessageEnt)
	if errs != nil {
		str := fmt.Sprintf("%v", errs)
		return PostMessageBadREquestResponse(str)
	}

	err = messageRepository.Create(userMessageEnt)
	if err != nil {
		return PostMessageInternalServerErrorResponse()
	}

	return PostMessageOKResponse(message)
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	// XXX このトークン認証が何回もありすぎてクソイケてない
	userTokenEnt, err := repositories.NewUserTokenRepository().GetByToken(p.Token)

	if err != nil {
		return GetMessagesInternalServerErrorResponse()
	}

	if userTokenEnt == nil {
		return GetMessageUnauthorizedResponse()
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
