package message

import (
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	message := repositories.NewUserMessageRepository()
	match := repositories.NewUserMatchRepository()

	// ユーザーID取得用
	ut, _ := t.GetByToken(p.Params.Token)
	// 自分が既にマッチングしている全てのUserIDを取得
	all, _ := match.FindAllByUserID(ut.UserID)

	// マッチしているユーザーかをきちんと確認する
	if CheckMatchUserID(all, p.UserID) {
		// メッセージの値の定義
		m := entities.UserMessage{
			UserID:    ut.UserID,
			PartnerID: p.UserID,
			Message:   p.Params.Message,
			CreatedAt: strfmt.DateTime(time.Now()),
			UpdatedAt: strfmt.DateTime(time.Now()),
		}
		// メッセージをインサートする
		err := message.Create(m)
		if err != nil {
			return si.NewPostLikeInternalServerError().WithPayload(
				&si.PostLikeInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	} else {
		return si.NewGetMessagesBadRequest().WithPayload(
			&si.GetMessagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
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
