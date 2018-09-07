package userlike

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

//	いいね！表示API
// 	GET {hostname}/api/1.0/likes
func getLikesInternalServerErrorResponse() middleware.Responder {
	return si.NewGetLikesInternalServerError().WithPayload(
		&si.GetLikesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getLikesUnauthorizedResponse() middleware.Responder {
	return si.NewGetLikesUnauthorized().WithPayload(
		&si.GetLikesUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func getLikesBadRequestResponse() middleware.Responder {
	return si.NewGetLikesBadRequest().WithPayload(
		&si.GetLikesBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

//	いいね！送信API
//	POST {hostname}/api/1.0/likes/{userID}
func postLikeInternalServerErrorResponse() middleware.Responder {
	return si.NewPostLikeInternalServerError().WithPayload(
		&si.PostLikeInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func postLikeUnauthorizedResponse() middleware.Responder {
	return si.NewPostLikeUnauthorized().WithPayload(
		&si.PostLikeUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func postLikeBadRequestResponses() middleware.Responder {
	return si.NewPostLikeBadRequest().WithPayload(
		&si.PostLikeBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

func postLikeOK() middleware.Responder {
	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:    "200",
			Message: "OK",
		})
}
