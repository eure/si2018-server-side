package message

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"time"
)

//- メッセージ送信API
//- POST {hostname}/api/1.0/messages/{userID}
//- TokenのValidation処理を実装してください
//- マッチングしていないとメッセージはできません
func PostMessage(p si.PostMessageParams) middleware.Responder {
	repMessage := repositories.NewUserMessageRepository()
	repUserToken := repositories.NewUserTokenRepository() // tokenからidを取得する
	repUserMatch := repositories.NewUserMatchRepository() //

	err := repUserToken.ValidateToken(p.Params.Token)
	if err != nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code: "401",
				Message: "Token Is Invalid",
			})
	}

	// Bad Request
	if p.UserID < 0 {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	if p.Params.Message == "" {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	var message entities.UserMessage

	// トークンからログインユーザーを取得
	loginUser, _ := repUserToken.GetByToken(p.Params.Token)

	// メッセージ送信相手とマッチングしていない場合, Bad Requestを返す
	userMatch, err := repUserMatch.Get(loginUser.UserID, p.UserID)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userMatch == nil{
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	message.UserID = loginUser.UserID
	message.PartnerID = p.UserID
	message.Message = p.Params.Message
	message.CreatedAt = strfmt.DateTime(time.Now())
	message.UpdatedAt = message.CreatedAt

	err = repMessage.Create(message)
	if err != nil{
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

//- メッセージ内容取得API
//- GET {hostname}/api/1.0/messages/{userID}
//- TokenのValidation処理を実装してください
//- ページネーションを実装してください
func GetMessages(p si.GetMessagesParams) middleware.Responder {
	repUserMessage := repositories.NewUserMessageRepository()
	repUserToken := repositories.NewUserTokenRepository()

	loginUser, _:= repUserToken.GetByToken(p.Token)
	
	userMessages, _:= repUserMessage.GetMessages(loginUser.UserID, p.UserID, int(*p.Limit), p.Latest, p.Oldest)

	if userMessages == nil {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Message is nothing",
			})
	}

	messages := entities.UserMessages(userMessages)

	messagesModel := messages.Build()

	return si.NewGetMessagesOK().WithPayload(messagesModel)
}
