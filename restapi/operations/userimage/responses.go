package userimage

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func postImageInternalServerErrorResponse() middleware.Responder {
	return si.NewPostImagesInternalServerError().WithPayload(
		&si.PostImagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func postImageUnauthorizedResponse() middleware.Responder {
	return si.NewPostImagesUnauthorized().WithPayload(
		&si.PostImagesUnauthorizedBody{
			Code:    "401",
			Message: "Tokan Is Invalid",
		})
}

func postImageBadRequestResponse(message string) middleware.Responder {
	return si.NewPostImagesBadRequest().WithPayload(
		&si.PostImagesBadRequestBody{
			Code:    "400",
			Message: message,
		})
}
