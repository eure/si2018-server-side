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
	//u := repositories.NewUserRepository()
	m := repositories.NewUserMatchRepository()
	
	// 作成するメッセージの宣言
	var message entities.UserMessage
	
	// Tokenの取得, 送信元UserIDの宣言
	FromUserToken := p.Params.Token
	FromUser, err := t.GetByToken(FromUserToken)
	if err != nil {
		return outPutPostStatus(500)
	}
	FromUserID := FromUser.UserID
	
	// メッセージ送信元UserIDの宣言
	ToUserID := p.UserID
	
	// メッセージの作成
	message.UserID = ToUserID
	message.PartnerID = FromUserID
	message.Message = p.Params.Message
	message.CreatedAt = strfmt.DateTime(time.Now())
	message.UpdatedAt = strfmt.DateTime(time.Now())
	
	// マッチ済みかどうか
	userMatch , err := m.Get(FromUserID,ToUserID)
	if userMatch != nil {
		r.Create(message)
	}
	
	userMatch , err = m.Get(ToUserID, FromUserID)
	if userMatch != nil {
		r.Create(message)
	}
	
	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {

	// UserIDと，PartnerIDがほしい
	// matchした相手を特定したのち，上記二つの情報プラスαをGetMeesagesに投げると良い

	// おなじみ TokenからUserIDを引っ張ってくる関数
	userR := repositories.NewUserRepository()
	user , _ := userR.GetByToken(p.Token)
	userid := user.ID

	// UserIDから，PartnerIDたちを取得
	userM := repositories.NewUserMatchRepository()
	var matchedusers []int64
	matchedusers , _ = userM.FindAllByUserID(userid)

	r := repositories.NewUserMessageRepository()

	var messages entities.UserMessages
	for _,m := range matchedusers {
		messages1partner , _ := r.GetMessages(userid,m,int(*p.Limit),p.Latest,p.Oldest)
		for _,message := range messages1partner {
			messages = append(messages, message)
		}
	}

	sEnt := messages.Build()
	return si.NewGetMessagesOK().WithPayload(sEnt)
}

func outPutPostStatus (num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
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