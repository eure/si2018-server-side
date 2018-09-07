package userimage

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"strings"
	"github.com/go-openapi/strfmt"
	"strconv"
	"os"
	"path"
	"bytes"
	
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	tr := repositories.NewUserTokenRepository()
	ur := repositories.NewUserRepository()
	uir := repositories.NewUserImageRepository()
	
	token, err := tr.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	if token == nil {
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code: "401",
				Message: "Unauthorized",
			})
	}

	usr, err := ur.GetByUserID(token.UserID)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	if usr == nil {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code: "400",
				Message: "Bad Request(Can't get User)",
			})
	}

	buf := bytes.NewBuffer(p.Params.Image)
	_, format, err := image.DecodeConfig(buf)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	var extension string
	if format == "png" || format == "jpeg" || format == "gif" {
		extension = "." + format

	} else {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code: "400",
				Message: "invalid image format",
			})
	}

	
	fileNameElems := []string{strconv.Itoa(int(usr.ID)), "test"}
	fileName := strings.Join(fileNameElems, "_") + extension
	assetPath := os.Getenv("ASSETS_PATH")

	newImage := entities.NewUserImage(usr.ID, path.Join(assetPath, fileName))
	println(newImage.UserID)
	println(newImage.Path)
	print(newImage.CreatedAt.String())
	print(newImage.UpdatedAt.String())

	// ファイル作成
	imageFile, err := os.Create(newImage.Path)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error(file create)",
			})
	}
	defer imageFile.Close()
	imageFile.Write(p.Params.Image)

	err = uir.Update(newImage)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	// image, err := uir.Update()

	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(newImage.Path),
		})
}
