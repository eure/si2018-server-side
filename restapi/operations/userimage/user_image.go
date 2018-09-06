package userimage

import (
	"bytes"
	"fmt"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
)










func PostImage(p si.PostImagesParams) middleware.Responder {
	/*
	1.	tokenのバリデーション
	2.	tokenから使用者のuseridを取得
	3.	画像をassetsの中に書き込む（本当はユニークな文字列にしなくてはいけない）
	4.	（画像を削除する）
	// userIDは送信者, partnerIDは受信者
	*/


	// Tokenがあるかどうか
	if p.Params.Token == "" {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "No Token",
			})
	}

	// tokenからuserIDを取得

	rToken := repositories.NewUserTokenRepository()
	entToken, errToken := rToken.GetByToken(p.Params.Token)
	if errToken != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entToken == nil {
		return si.NewPostMessageUnauthorized().WithPayload(
			&si.PostMessageUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized Token",
			})
	}

	sEntToken := entToken.Build()

	rImage := repositories.NewUserImageRepository()

	assets_path := os.Getenv("ASSETS_PATH")

	// 拡張子の判別
	/*
	バイナリの先頭部分からjpg, pngを判定する（それ以外の場合はbad requestを返す）
	jpegファイルの場合は\xff\xd8
	pngファイルの場合は\x89\x50\x4e\x47\x0d\x0a\x1a\x0a
	*/
	var extension string
	if bytes.HasPrefix(p.Params.Image,[]byte{0xff, 0xd8}) {
		extension = "jpg"
	} else if bytes.HasPrefix(p.Params.Image,[]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}) {
		extension = "png"
	} else {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				"400",
				"Bad Request",
			})
	}

	filename := strconv.FormatInt(sEntToken.UserID, 10)+"."+extension
	file, errOpen := os.OpenFile(assets_path+filename, os.O_WRONLY|os.O_CREATE, 0666)

	if errOpen != nil {
		fmt.Println(errOpen)
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	defer file.Close()

	_, errOpen = file.Write(p.Params.Image)

	if errOpen  != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}


	var userImage entities.UserImage
	userImage.UserID = sEntToken.UserID
	userImage.Path = "https://127.0.0.1:8888/assets/" + filename

	errUpdate := rImage.Update(userImage)

	if errUpdate != nil {
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			strfmt.URI("https://127.0.0.1:8888/assets/" + filename),
		})
}
