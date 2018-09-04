package message

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	return si.NewPostMessageOK()
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {

	// UserIDと，PartnerIDがほしい
	// matchした相手を特定したのち，上記二つの情報をGetMeesagesに投げると良い

	// おなじみ TokenからUserIDを引っ張ってくる関数
	userR := repositories.NewUserRepository()
	user , _ := userR.GetByToken(p.Token)
	userid := user.ID

	// UserIDから，PartnerIDたちを取得
	userM := repositories.NewUserMatchRepository()
	var matchedusers []int64
	matchedusers , _ = userM.FindAllByUserID(userid)

	r := repositories.NewUserMessageRepository()

	var matcheduser entities.UserMessage
	var messages []entities.UserMessages
	for _,m := range matchedusers {
		matcheduser = r.GetMessages(userid,m)
		messages = append (messages , matcheduser)
	}

	sEnt := messages.Build()

	return si.NewGetMessagesOK()
}
