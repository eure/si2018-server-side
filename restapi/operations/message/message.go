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
	t := repositories.NewUserTokenRepository()
	message := repositories.NewUserMessageRepository()
	match := repositories.NewUserMatchRepository()

	// ユーザーID取得用
	ut, _ := t.GetByToken(p.Token)
	// matchingしているユーザーの取得
	all, _ := match.FindAllByUserID(ut.UserID)

	// 明示的に宣言
	var m entities.UserMessages
	// マッチしているユーザーかをきちんと確認する
	if CheckMatchUserID(all, p.UserID) {
		m, _ = message.GetMessages(ut.UserID, p.UserID, int(*p.Limit), p.Latest, p.Oldest)
		sEnt := m.Build()
		return si.NewGetMessagesOK().WithPayload(sEnt)
	} else {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
}

//マッチングをしているか確認に使う関数
func CheckMatchUserID(matchID []int64, userID int64) bool {
	for _, m := range matchID {
		if m == userID {
			return true
		}
	}
	return false
}
