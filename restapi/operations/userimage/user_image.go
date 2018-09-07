package userimage

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"strconv"
	"time"
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
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Required",
			})
	}

	// tokenからuserIDを取得

	rToken := repositories.NewUserTokenRepository()
	entToken, errToken := rToken.GetByToken(p.Params.Token)
	if errToken != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entToken == nil {
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
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
	img := p.Params.Image
	// imgが5MB以上の時は処理しない。
	if len(img) > 1048576*5 {
		return si.NewPostImagesRequestEntityTooLarge().WithPayload(
			&si.PostImagesRequestEntityTooLargeBody{
				Code:    "413",
				Message: "Payload Too Large",
			})
	}

	var imgDecoded image.Image
	var errDecode error
	var extension string
	var errEncode error

	buf := new(bytes.Buffer)
	switch {
	case bytes.HasPrefix(img, []byte{0xff, 0xd8}):
		extension = "jpg"
		imgDecoded, errDecode = jpeg.Decode(bytes.NewReader(img))
	case bytes.HasPrefix(img, []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}):
		extension = "png"
		imgDecoded, errDecode = png.Decode(bytes.NewReader(img))
	default:
		return si.NewPostImagesUnsupportedMediaType().WithPayload(
			&si.PostImagesUnsupportedMediaTypeBody{
				Code:    "415",
				Message: "Unsupported Media Type",
			})
	}

	if errDecode != nil || imgDecoded == nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	imgResizedDecoded := resize.Thumbnail(400, 400, imgDecoded, resize.Lanczos3)

	switch extension {
	case "jpg":
		errEncode = jpeg.Encode(buf, imgResizedDecoded, nil)
	case "png":
		errEncode = png.Encode(buf, imgResizedDecoded)
	}

	if errEncode != nil || buf == nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	img = buf.Bytes()

	// ファイル名……userID+時刻をmd5hashに変換
	filestring := strconv.FormatInt(sEntToken.UserID, 10) + "_" + time.Now().String()
	filehash := md5.Sum([]byte(filestring))
	filename := hex.EncodeToString(filehash[:]) + "." + extension

	file, errOpen := os.OpenFile(assets_path+filename, os.O_WRONLY|os.O_CREATE, 0666)

	if errOpen != nil {
		fmt.Println(errOpen)
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	defer file.Close()

	_, errOpen = file.Write(img)

	if errOpen != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	var userImage entities.UserImage
	baseAddress := "https://127.0.0.1:8888/assets/"
	userImage.UserID = sEntToken.UserID
	userImage.Path = baseAddress + filename

	errUpdate := rImage.Update(userImage)

	if errUpdate != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(baseAddress + filename),
		})
}
