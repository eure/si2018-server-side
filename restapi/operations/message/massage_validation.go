package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

//	メッセージ送信API
// 	POST {hostname}/api/1.0/messages/{userID}
func postMessageInternalServerErrorResponse() middleware.Responder {
	return si.NewPostMessageInternalServerError().WithPayload(
		&si.PostMessageInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func postMessageUnauthorizedResponse() middleware.Responder {
	return si.NewPostMessageUnauthorized().WithPayload(
		&si.PostMessageUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func postMessageBadRequestResponses() middleware.Responder {
	return si.NewPostMessageBadRequest().WithPayload(
		&si.PostMessageBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

// メッセージは重複していますという400番エラー
func postMessageDuplicatedRequestResponses() middleware.Responder {
	return si.NewPostMessageBadRequest().WithPayload(
		&si.PostMessageBadRequestBody{
			Code:    "400",
			Message: "Message is duplicated",
		})
}

// 	メッセージ内容取得API
//	GET {hostname}/api/1.0/messages/{userID}
func getMessagesInternalServerErrorResponse() middleware.Responder {
	return si.NewGetMessagesInternalServerError().WithPayload(
		&si.GetMessagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getMessagesUnauthorizedResponse() middleware.Responder {
	return si.NewGetMessagesUnauthorized().WithPayload(
		&si.GetMessagesUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func getMessagesBadRequestResponse() middleware.Responder {
	return si.NewGetMessagesBadRequest().WithPayload(
		&si.GetMessagesBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}
