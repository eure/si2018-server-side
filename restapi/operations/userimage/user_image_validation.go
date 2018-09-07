package userimage

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// 	プロフィール写真の更新
//	POST {hostname}/api/1.0/images
func postImagesInternalServerErrorResponse() middleware.Responder {
	return si.NewPostImagesInternalServerError().WithPayload(
		&si.PostImagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func postImagesUnauthorizedResponse() middleware.Responder {
	return si.NewPostImagesUnauthorized().WithPayload(
		&si.PostImagesUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

func postImagesBadRequestResponse() middleware.Responder {
	return si.NewPostImagesBadRequest().WithPayload(
		&si.PostImagesBadRequestBody{
			Code:    "401",
			Message: "Bad Request",
		})
}

func postImagesOKResponse(imagePath string) middleware.Responder {
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(imagePath),
		})
}
