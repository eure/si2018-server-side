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

	// Message
	if p.Params.Message == "" {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// トークンからログインユーザーを取得
	loginUser, err := repUserToken.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// メッセージ送信相手とマッチングしているか確認
	userMatch, err := repUserMatch.Get(loginUser.UserID, p.UserID)
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	// メッセージ送信相手とマッチングしていない場合, Bad Requestを返す
	if userMatch == nil{
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// UserMessageの作成
	var message entities.UserMessage
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

	err := repUserToken.ValidateToken(p.Token)
	if err != nil {
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code: "401",
				Message: "Token Is Invalid",
			})
	}

	// limitが0の場合、メッセージを全取得してしまう => []を返す
	if p.Limit != nil {
		if *p.Limit == 0 {
			return si.NewGetMessagesOK().WithPayload(nil)
		} else if *p.Limit < 0 {
			return si.NewGetMessagesBadRequest().WithPayload(
				&si.GetMessagesBadRequestBody{
					Code:    "400",
					Message: "Bad Request",
				})
		}
	}

	// OldestとLatestが存在する場合
	if (p.Oldest != nil) && (p.Latest != nil) {
		// LatestがOldestより時間的に後だったら　
		if time.Time(*p.Latest).After(time.Time(*p.Oldest)) {
			return si.NewGetMessagesBadRequest().WithPayload(
				&si.GetMessagesBadRequestBody{
					Code: "400",
					Message: "Bad Request",
				})
		}

	}

	// tokenからユーザーの取得
	loginUser, err:= repUserToken.GetByToken(p.Token)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// メッセージの取得
	userMessages, err:= repUserMessage.GetMessages(loginUser.UserID, p.UserID, int(*p.Limit), p.Latest, p.Oldest)
	if err != nil {
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// UserMessageの配列をUserMessagesにキャスト
	messages := entities.UserMessages(userMessages)

	messagesModel := messages.Build()

	return si.NewGetMessagesOK().WithPayload(messagesModel)
}
