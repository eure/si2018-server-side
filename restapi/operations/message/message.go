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
	// userIDは送信者, partnerIDは受信者
	*/


	// Tokenがあるかどうか
	if p.Params.Token == "" {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "No Token",
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
				Message: "Unauthorized Token",
			})
	}

	sEntToken := entToken.Build()

	rMatch := repositories.NewUserMatchRepository()
	ids, err := rMatch.FindAllByUserID(sEntToken.UserID)

	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
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
	if id_contains {
		rMessage := repositories.NewUserMessageRepository()
		var message entities.UserMessage
		message.UserID = sEntToken.UserID
		message.PartnerID = p.UserID
		message.Message = p.Params.Message
		message.CreatedAt = strfmt.DateTime(time.Now())
		message.UpdatedAt = message.CreatedAt
		rMessage.Create(message)

		return si.NewPostMessageOK().WithPayload(
			&si.PostMessageOKBody{
				"200",
				"OK",
			})
	} else {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})

	}

}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	/*
	1. 	Tokenのバリデーション
	2.	Tokenから使用者のuserIDを取得
	3.	メッセージ送り先のuseridを取得する
	badrequestの必要がない？
	*/

	// Tokenがあるかどうか
	if p.Token == "" {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "No Token",
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
				Message: "Unauthorized Token",
			})
	}

	sEntToken := entToken.Build()

	rMessage := repositories.NewUserMessageRepository()

	limit := int(*(p.Limit))
	latest := time.Time(*(p.Latest))
	oldest := time.Time(*(p.Oldest))
	now := time.Now()
	// oldest latest now ->（時系列）の順になっていないとおかしい
	if limit < 0 || !(oldest.Before(latest) && latest.Before(now)) {
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
