package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/runtime/middleware"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	return si.NewPostMessageOK()
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	message := repositories.NewUserMessageRepository()
	match := repositories.NewUserMatchRepository()

	ut, _ := t.GetByToken(p.Token)
	all, _ := match.FindAllByUserID(ut.UserID)

  var m entities.UserMessages
	if CheckLikeUserID(all, p.UserID){
		m, _ = message.GetMessages(ut.UserID, p.UserID, int(*p.Limit), p.Latest, p.Oldest)
		sEnt := m.Build()
		return si.NewGetMessagesOK().WithPayload(sEnt)
	}else{
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
}

//マッチングをしているか確認に使う関数
func CheckLikeUserID(likeID []int64, userID int64) bool{
	for _, m := range likeID{
		if m == userID{
			return true
		}
	}
	return false
}
