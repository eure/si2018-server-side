package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/strfmt"
	"time"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {

	// レポジトリを初期化する
	tokenR := repositories.NewUserTokenRepository()
	userMessageR := repositories.NewUserMessageRepository()

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)

	// 401エラー
	if tokenEnt == nil {
		si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code: "401",
				Message: "Token is invalid",
			})
	}

	// 500エラー
	if err != nil {
		si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	// メッセージ作成用の構造体を作成
	userMessageEnt := entities.UserMessage{
		UserID: tokenEnt.UserID,
		PartnerID: p.UserID,
		Message: p.Params.Message,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}

	// メッセージを作成する
	err = userMessageR.Create(userMessageEnt)

	// 500エラー
	if err != nil {
		si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	// 結果を返す
	return si.NewPostMessageOK().WithPayload(
		&si.PostMessageOKBody{
			Code: "200",
			Message: "OK",
		})
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {

	// レポジトリを初期化する
	tokenR := repositories.NewUserTokenRepository()
	usermessageR := repositories.NewUserMessageRepository()

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Token)

	// 401エラー
	if tokenEnt == nil {
		si.NewGetMessagesUnauthorized().WithPayload(
			&si.GetMessagesUnauthorizedBody{
				Code: "401",
				Message: "Token is invalid",
			})
	}

	// 500エラー
	if err != nil {
		si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	// limitのデフォルト値を100に設定する
	var limit int
	if &p.Limit == nil {
		limit = 100
	} else {
		limit = int(*p.Limit)
	}

	// messageの取得
	messageEnts, err := usermessageR.GetMessages(tokenEnt.UserID, p.UserID, limit, p.Latest, p.Oldest)

	// 500エラー
	if err != nil {
		si.NewGetMessagesInternalServerError().WithPayload(
			&si.GetMessagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	var messageEntities entities.UserMessages

	messageEntities = messageEnts

	// モデルに変換
	messages := messageEntities.Build()

	// 結果を返す
	return si.NewGetMessagesOK().WithPayload(messages)
}
