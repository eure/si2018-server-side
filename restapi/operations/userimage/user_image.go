package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"os"
	"strconv"
	"github.com/go-openapi/strfmt"
	"github.com/eure/si2018-server-side/entities"
)

const baseUri = "http://localhhost:8888/"

func PostImage(p si.PostImagesParams) middleware.Responder {
	// Tokenのユーザが存在しない -> 401
	tokenR        := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	if tokenEnt == nil{
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code   : "401",
				Message: "Token Is Invalid",
			})
	}

	// DONE: Base64から画像ファイルを作成
	// TODO: ファイル名をランダムな文字列にしたい
	fileName := strconv.Itoa(int(tokenEnt.UserID))+".jpg"
	filePath := "assets/"+fileName
	file, err := os.Create(filePath)
	if err != nil {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	defer file.Close()
	file.Write(p.Params.Image)

	// DONE: DBの更新
	// image_uriを返す
	uri := baseUri + filePath
	imageR := repositories.NewUserImageRepository()
	tmp := entities.UserImage{
		UserID: tokenEnt.UserID,
		Path: uri,
	}
	err = imageR.Update(tmp)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	imageEnt, err := imageR.GetByUserID(tokenEnt.UserID)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// TODO: 正常に終了したら
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI:strfmt.URI(imageEnt.Path),
		})
}
