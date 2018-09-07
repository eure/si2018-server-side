package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	UserTokenRepository := repositories.NewUserTokenRepository()
	UserImageRepository := repositories.NewUserImageRepository()
	assetsPath := os.Getenv("ASSETS_PATH")
	UserByToken, err := UserTokenRepository.GetByToken(p.Params.Token)
	if err != nil {
		si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if UserByToken == nil {
		si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	ImgName := strconv.Itoa(int(UserByToken.UserID)) + "_image.png"
	ImgPath := assetsPath + ImgName
	file, err := os.Create(ImgPath)
	if err != nil {
		si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	defer file.Close()
	file.Write(p.Params.Image)
	UserImage, err := UserImageRepository.GetByUserID(UserByToken.UserID)
	if err != nil {
		si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if UserImage == nil {
		si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	UserImage.Path = ImgPath
	err = UserImageRepository.Update(*UserImage)
	if err != nil {
		si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(ImgPath),
		})
}
