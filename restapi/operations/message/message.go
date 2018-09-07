package message

import (
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// マッチングしている指定のユーザにメッセージ送信
func PostMessage(p si.PostMessageParams) middleware.Responder {
	UserTokenRepository := repositories.NewUserTokenRepository()
	UserMessageRepository := repositories.NewUserMessageRepository()
	UserMatchRepository := repositories.NewUserMatchRepository()

	Token := p.Params.Token
	UserByToken, err := UserTokenRepository.GetByToken(Token)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	// GetByTokenで返り値があれば正しいトークンである
	if UserByToken == nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	// 文字数バリデーション
	if len(p.Params.Message) >= 1000 {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	UserID := UserByToken.UserID
	PartnerID := p.UserID

	if PartnerID <= 0 {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	//
	Match, err := UserMatchRepository.Get(UserID, PartnerID)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	var InsertMessage entities.UserMessage
	if Match != nil {
		InsertMessage = entities.UserMessage{
			UserID:    UserID,
			PartnerID: PartnerID,
			Message:   p.Params.Message,
			CreatedAt: strfmt.DateTime(time.Now()),
			UpdatedAt: strfmt.DateTime(time.Now()),
		}
	}
	err = UserMessageRepository.Create(InsertMessage)
	if err != nil {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	UserMessageRepository := repositories.NewUserMessageRepository()
	UserTokenRepository := repositories.NewUserTokenRepository()

	UserByToken, err := UserTokenRepository.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if UserByToken == nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	if int(*(p.Limit)) <= 0 {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	UserID := UserByToken.UserID
	PartnerID := p.UserID

	Messages, err := UserMessageRepository.GetMessages(UserID, PartnerID, int(*(p.Limit)), p.Latest, p.Oldest)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	// MessageがないとBadになってしまう
	/*
		if Messages == nil {
			return si.NewGetMessagesBadRequest().WithPayload(
				&si.GetMessagesBadRequestBody{
					Code:    "400",
					Message: "Bad Request",
				})
		}
	*/

	UserMessages := entities.UserMessages(Messages)
	BuildedUserMessages := UserMessages.Build()
	return si.NewGetMessagesOK().WithPayload(BuildedUserMessages)
}
