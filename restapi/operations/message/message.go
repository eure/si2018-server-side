package message

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/runtime/middleware"
	
	"time"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	r := repositories.NewUserMessageRepository()
	t := repositories.NewUserTokenRepository()
	m := repositories.NewUserMatchRepository()
	
	// fromUser == loginUser
	// tokenから UserToken entitiesの取得 (Validation)
	token := p.Params.Token
	fromUser, err := t.GetByToken(token)
	if err != nil {
		return outPutPostStatus(500)
	}
	fromUserID := fromUser.UserID
	
	// メッセージ送信元UserIDの宣言
	toUserID := p.UserID
	
	// メッセージの作成
	var message entities.UserMessage
	message.UserID = toUserID
	message.PartnerID = fromUserID
	message.Message = p.Params.Message
	message.CreatedAt = strfmt.DateTime(time.Now())
	message.UpdatedAt = strfmt.DateTime(time.Now())
	
	// マッチ済みかどうか
	userMatch , err := m.Get(fromUserID,toUserID)
	if err != nil {
		return outPutPostStatus(500)
	}
	if userMatch != nil {
		r.Create(message)
	}
	
	userMatch , err = m.Get(toUserID, fromUserID)
	if err != nil {
		return outPutPostStatus(500)
	}
	if userMatch != nil {
		r.Create(message)
	} else {
		return outPutPostStatus(400)
	}
	
	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: message.Message,
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	r := repositories.NewUserMessageRepository()
	m := repositories.NewUserMatchRepository()
	
	// loginUserのUserToken entitiesを取得(Validation)
	loginUser , err := t.GetByToken(p.Token)
	if err != nil {
		outPutGetStatus(500)
	}
	if loginUser == nil {
		outPutGetStatus(401)
	}
	
	// UserIDから，PartnerIDたちを取得
	loginUserID := loginUser.UserID
	var matchedUserIDs []int64
	matchedUserIDs , err = m.FindAllByUserID(loginUserID)
	if err != nil {
		return outPutGetStatus(500)
	}

	var messages entities.UserMessages
	for _,m := range matchedUserIDs {
		messages1partner , err := r.GetMessages(loginUserID,m,int(*p.Limit),p.Latest,p.Oldest)
		if err != nil {
			return outPutGetStatus(500)
		}
		if messages1partner == nil {
			return outPutGetStatus(400)
		}
		
		for _,message := range messages1partner {
			messages = append(messages, message)
		}
	}

	sEnt := messages.Build()
	return si.NewGetMessagesOK().WithPayload(sEnt)
}

func outPutGetStatus (num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	case 401:
		return si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	case 400:
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	return nil
}

func outPutPostStatus (num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	case 401:
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗)",
			})
	case 400:
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	return nil
}
