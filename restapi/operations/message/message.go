package message

import (
	"fmt"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"time"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	/*
		1.	tokenのバリデーション
		2.	tokenから使用者のuseridを取得
		3.	メッセージ送り先のuseridがマッチング済みの相手かどうかを確認する
		4.	メッセージを送信する
		// userIDは送信元, partnerIDは送信先
	*/

	// Tokenがあるかどうか
	if p.Params.Token == "" {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Required",
			})
	}

	// tokenからuserIDを取得

	rToken := repositories.NewUserTokenRepository()
	entToken, errToken := rToken.GetByToken(p.Params.Token)
	if errToken != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entToken == nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	sEntToken := entToken.Build()

	// userIDがマッチしているIDを見つける
	rMatch := repositories.NewUserMatchRepository()
	ids, errMatch := rMatch.FindAllByUserID(sEntToken.UserID)

	if errMatch != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	var id_contains bool
	for _, id := range ids {
		if id == p.UserID {
			id_contains = true
			break
		}
	}

	if !id_contains {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	rMessage := repositories.NewUserMessageRepository()
	var message entities.UserMessage
	message.UserID = sEntToken.UserID
	message.PartnerID = p.UserID
	message.Message = p.Params.Message
	message.CreatedAt = strfmt.DateTime(time.Now())
	message.UpdatedAt = message.CreatedAt
	errorCreate := rMessage.Create(message)

	if errorCreate != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			"200",
			"OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	/*
		1. 	Tokenのバリデーション
		2.	Tokenから使用者のuserIDを取得
		3.	メッセージ送り先のuseridがマッチング済みの相手かどうかを確認する
		4.	メッセージの送信
	*/

	// Tokenがあるかどうか
	if p.Token == "" {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Required",
			})
	}

	// tokenからuserIDを取得

	rToken := repositories.NewUserTokenRepository()
	entToken, errToken := rToken.GetByToken(p.Token)
	if errToken != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entToken == nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	sEntToken := entToken.Build()

	// userIDがマッチしているIDを見つける
	rMatch := repositories.NewUserMatchRepository()
	ids, err := rMatch.FindAllByUserID(sEntToken.UserID)

	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	fmt.Println(ids)
	fmt.Println(p.UserID)
	var id_contains bool
	for _, id := range ids {
		if id == p.UserID {
			id_contains = true
			break
		}
	}

	if !id_contains {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	rMessage := repositories.NewUserMessageRepository()

	limit := int(*(p.Limit))

	if limit < 0 {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	entMessage, errMessage := rMessage.GetMessages(sEntToken.UserID, p.UserID, limit, p.Latest, p.Oldest)

	if errMessage != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	userMessages := entities.UserMessages(entMessage)

	sMsgs := userMessages.Build()

	return si.NewGetMessagesOK().WithPayload(sMsgs)

}
