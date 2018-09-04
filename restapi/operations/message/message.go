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
	repUserMatch := repositories.NewUserMatchRepository()

	var message entities.UserMessage

	loginUser, _ := repUserToken.GetByToken(p.Params.Token)

	//マッチングしていない場合エラーを返す
	userMatch, err := repUserMatch.Get(loginUser.UserID, p.UserID)
	if userMatch == nil{
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Can't sent message for not matched user",
			})
	}
	if err != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
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
			Message: "Posted Your Message",
		})
}


//- メッセージ内容取得API
//- GET {hostname}/api/1.0/messages/{userID}
//- TokenのValidation処理を実装してください
//- ページネーションを実装してください
func GetMessages(p si.GetMessagesParams) middleware.Responder {
	return si.NewGetMessagesOK()
}
