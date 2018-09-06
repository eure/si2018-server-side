package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"os"
	"github.com/go-openapi/strfmt"
	"github.com/eure/si2018-server-side/entities"
	"encoding/hex"
	"strings"
	"math/rand"
	"time"
)

const baseUri = "http://localhhost:8888/"

func PostImage(p si.PostImagesParams) middleware.Responder {
	tokenR        := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// Tokenのユーザが存在しない -> 401
	if tokenEnt == nil{
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code   : "401",
				Message: "Token Is Invalid",
			})
	}
	// ランダムな文字列作成
	rand.Seed(time.Now().UnixNano())
	n := 10
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	randString := string(b)
	// Base64から画像ファイルを作成
	fileString := hex.EncodeToString(p.Params.Image)
	var fileName string
	// ファイルの形式を確認
	if strings.HasPrefix(fileString, "ffd8") {
		fileName = randString+".jpg"
	}else if strings.HasPrefix(fileString, "89504e47"){
		fileName = randString+".png"
	}else if strings.HasPrefix(fileString, "47494638"){
		fileName = randString+".gif"
	}
	filePath := "assets/"+fileName
	// ファイル作成
	file, err := os.Create(filePath)
	if err != nil {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	defer file.Close()
	_, err = file.Write(p.Params.Image)
	if err != nil {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	// DBの更新
	imageUri := baseUri + filePath
	imageR := repositories.NewUserImageRepository()
	tmp := entities.UserImage{
		UserID: tokenEnt.UserID,
		Path  : imageUri,
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
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code   : "500",
				Message: "Internal Server Error",
			})
	}
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI:strfmt.URI(imageEnt.Path),
		})
}
