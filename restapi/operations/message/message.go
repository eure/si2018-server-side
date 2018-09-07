package message

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

var (
	messageRepo = repositories.NewUserMessageRepository()
	tokenRepo   = repositories.NewUserTokenRepository()
	matchRepo   = repositories.NewUserMatchRepository()
)

//	メッセージ送信API
func PostMessage(p si.PostMessageParams) middleware.Responder {
	// 入力値のValidation処理をします。
	partnerID := p.UserID
	if partnerID <= 0 {
		return postMessageBadRequestResponses()
	}

	message := p.Params.Message
	if message == "" {
		return postMessageBadRequestResponses()
	}

	token := p.Params.Token

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return postMessageInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return postMessageUnauthorizedResponse()
	}

	id := tokenOwner.UserID

	// 連続投稿をさせたくないため、過去にあった内容と入力された内容が同じなら処理を止める。
	pastMessage, err := messageRepo.GetLastMessages(id)
	if err != nil {
		return postMessageInternalServerErrorResponse()
	} else if pastMessage.Message == message {
		return postMessageDuplicatedRequestResponses()
	}

	// ユーザーとお相手がマッチングしているか検証します。
	match, err := matchRepo.Get(id, partnerID)
	if err != nil {
		return postMessageInternalServerErrorResponse()
	} else if match == nil {
		return postMessageBadRequestResponses()
	}

	// メッセージを送信します。
	addMessage := entities.UserMessage{
		UserID:    id,
		PartnerID: partnerID,
		Message:   message,
	}
	err = messageRepo.Create(addMessage)
	if err != nil {
		return postMessageInternalServerErrorResponse()
	}

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "Posted Your Message",
		})
}

// 	メッセージ内容取得API
func GetMessages(p si.GetMessagesParams) middleware.Responder {
	// 入力値のValidation処理をします。
	partnerID := p.UserID
	if partnerID <= 0 {
		return getMessagesBadRequestResponse()
	}

	latest := p.Latest
	oldest := p.Oldest
	limit := int(*p.Limit)
	token := p.Token

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return getMessagesInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return getMessagesUnauthorizedResponse()
	}

	id := tokenOwner.UserID

	// ユーザーがお相手とやりとりしたメッセージを取得します
	messages, err := messageRepo.GetMessages(id, partnerID, limit, latest, oldest)
	if err != nil {
		return getMessagesInternalServerErrorResponse()
	}

	var ents entities.UserMessages
	ents = messages
	sEnt := ents.Build()
	return si.NewGetMessagesOK().WithPayload(sEnt)
}
