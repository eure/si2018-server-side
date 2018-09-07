package message

import (
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

// DB アクセス: 3 回
// 計算量: O(1)
func PostMessage(p si.PostMessageParams) middleware.Responder {
	var err error
	if p.Params.Token == "" {
		return si.PostMessageThrowBadRequest("トークンが与えられていません")
	}
	if p.Params.Message == "" {
		return si.PostMessageThrowBadRequest("メッセージが空です")
	}
	if len(p.Params.Message) > 3000 {
		return si.PostMessageThrowBadRequest("メッセージが長すぎます")
	}
	messageRepo := repositories.NewUserMessageRepository()
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		t, err := tokenRepo.GetByToken(p.Params.Token)
		// トークン認証
		if err != nil {
			return si.PostMessageThrowInternalServerError(err)
		}
		if t == nil {
			return si.PostMessageThrowUnauthorized()
		}
		id = t.UserID
	}
	// マッチしているかの確認
	{
		matchRepo := repositories.NewUserMatchRepository()
		match, err := matchRepo.Get(p.UserID, id)
		if err != nil {
			return si.PostMessageThrowInternalServerError(err)
		}
		if match == nil {
			return si.PostMessageThrowBadRequest("マッチしていない相手にはメッセージを送ることはできません")
		}
	}
	// メッセージを書きこむ
	now := strfmt.DateTime(time.Now())
	mes := entities.UserMessage{
		UserID:    id,
		PartnerID: p.UserID,
		Message:   p.Params.Message,
		UpdatedAt: now,
		CreatedAt: now}
	err = messageRepo.Create(mes)
	if err != nil {
		return si.PostMessageThrowInternalServerError(err)
	}
	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code:    "200",
			Message: "OK",
		})
}

// DB アクセス: 2 回
// 計算量: O(1)
func GetMessages(p si.GetMessagesParams) middleware.Responder {
	var err error
	messageRepo := repositories.NewUserMessageRepository()
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Token)
		if err != nil {
			return si.GetMessagesThrowInternalServerError(err)
		}
		if token == nil {
			return si.GetMessagesThrowUnauthorized()
		}
		id = token.UserID
	}
	// p.Limit のデフォルトは 100 (restapi/summerintern/get_messages_parameters.go)
	var limit int
	if p.Limit == nil {
		limit = 100
	} else {
		limit = int(*p.Limit)
	}
	// メッセージの取得
	message, err := messageRepo.GetMessages(p.UserID, id, limit, p.Latest, p.Oldest)
	if err != nil {
		return si.GetMessagesThrowInternalServerError(err)
	}
	ent := entities.UserMessages(message)
	return si.NewGetMessagesOK().WithPayload(ent.Build())
}
