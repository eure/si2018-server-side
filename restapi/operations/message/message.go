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
