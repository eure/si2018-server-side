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
		return si.PostMessageThrowBadRequest("missing token")
	}
	if p.Params.Message == "" {
		return si.PostMessageThrowBadRequest("empty message")
	}
	if len(p.Params.Message) > 3000 {
		return si.PostMessageThrowBadRequest("message too long")
	}
	messageRepo := repositories.NewUserMessageRepository()
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		t, err := tokenRepo.GetByToken(p.Params.Token)
		// トークン認証
		if err != nil {
			return si.PostMessageThrowInternalServerError("GetByToken", err)
		}
		if t == nil {
			return si.PostMessageThrowUnauthorized("GetByToken failed")
		}
		id = t.UserID
	}
	// マッチしているかの確認
	{
		matchRepo := repositories.NewUserMatchRepository()
		match, err := matchRepo.Get(p.UserID, id)
		if err != nil {
			return si.PostMessageThrowInternalServerError("Get", err)
		}
		if match == nil {
			return si.PostMessageThrowBadRequest("Get failed")
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
		return si.PostMessageThrowInternalServerError("Create", err)
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
			return si.GetMessagesThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return si.GetMessagesThrowUnauthorized("GetByToken failed")
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
		return si.GetMessagesThrowInternalServerError("GetMessages", err)
	}
	ent := entities.UserMessages(message)
	return si.NewGetMessagesOK().WithPayload(ent.Build())
}
