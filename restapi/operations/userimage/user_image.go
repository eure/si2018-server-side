package userimage

import (
	"bytes"
	"fmt"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/nfnt/resize"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
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
	400*400におさまるようにリサイズします。
	*/
	var extension string
	img := p.Params.Image
	// imgが5MB以上の時は処理しない。413に対応するメソッドがない……。
	if len(img) > 1048576*5 {
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				"413",
				"Payload Too Large",
			})
	}

	// jpegの場合
	if bytes.HasPrefix(img,[]byte{0xff, 0xd8}) {
		extension = "jpg"
		imgDecoded, errJpg := jpeg.Decode(bytes.NewReader(img))
		if errJpg != nil {
			return si.NewPostMessageInternalServerError().WithPayload(
				&si.PostMessageInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
		imgResizedDecoded := resize.Thumbnail(400,400,imgDecoded, resize.Lanczos3)
		buf := new(bytes.Buffer)
		errEncode := jpeg.Encode(buf, imgResizedDecoded, nil)
		if errEncode != nil {
			return si.NewPostMessageInternalServerError().WithPayload(
				&si.PostMessageInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
		img = buf.Bytes()
	// pngの場合
	} else if bytes.HasPrefix(img,[]byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}) {
		extension = "png"
		imgDecoded, errPng := png.Decode(bytes.NewReader(img))
		if errPng != nil {
			return si.NewPostMessageInternalServerError().WithPayload(
				&si.PostMessageInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
		imgResizedDecoded := resize.Thumbnail(400,400,imgDecoded, resize.Lanczos3)
		buf := new(bytes.Buffer)
		errEncode := png.Encode(buf, imgResizedDecoded)
		if errEncode != nil {
			return si.NewPostMessageInternalServerError().WithPayload(
				&si.PostMessageInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
		img = buf.Bytes()
	} else { // imgがjpg, png以外の時は処理しない。415に対応するメソッドがない……。
		return si.NewPostMessageBadRequest().WithPayload(
			&si.PostMessageBadRequestBody{
				"415",
				"Unsupported Media Type",
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

	_, errOpen = file.Write(img)

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
