package userimage

import (
	"fmt"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	userTokenEnt, err := repositories.NewUserTokenRepository().GetByToken(p.Params.Token)
	if err != nil {
		return postImageInternalServerErrorResponse()
	}
	if userTokenEnt == nil {
		return postImageUnauthorizedResponse()
	}

	userID := userTokenEnt.UserID

	userImageRepository := repositories.NewUserImageRepository()
	typeName, err := userImageRepository.ImageValidation(p.Params.Image)
	if err != nil {
		str := fmt.Sprintf("%v", err)
		return postImageBadRequestResponse(str)
	}

	imagePath, err := userImageRepository.SaveImageAssets(p.Params.Image, userID, typeName)
	if err != nil {
		return postImageInternalServerErrorResponse()
	}

	userImage, err := userImageRepository.GetByUserID(userID)
	if err != nil {
		return postImageInternalServerErrorResponse()
	}

	imagePath = "https://si-2018-011.eure.jp/" + imagePath
	userImage.Path = imagePath
	err = userImageRepository.Update(*userImage)
	if err != nil {
		return postImageInternalServerErrorResponse()
	}

	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(userImage.Path),
		})
}

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
