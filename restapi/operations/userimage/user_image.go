package userimage

import (
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
	3.	画像を送信する
	（4.	画像を削除する）
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

	// ファイルを書き込み用にオープン

	// 拡張子判別
	// base64をデコードしてbyte
	// byteにしてからio.readerの形にする
	// image.Decodeconfigでformat判別
	/*
	f, _ := os.Open("path/to/image")

	defer f.Close()

	*/

	filename := strconv.FormatInt(sEntToken.UserID, 10)+".jpg"
	file, err := os.OpenFile(assets_path+filename, os.O_WRONLY|os.O_CREATE, 0666)

	/*
	byteData, errDecode := base64.StdEncoding.DecodeString(p.Params.Image)
	_, format, err := image.DecodeConfig(strings.NewReader())
	*/
	// pngファイルの場合はformatがpng
	// jpegファイルの場合はformatがjpg
	// gifファイルの場合はformatがgif
	/*
	if err != nil {
		// 画像フォーマットではない場合はエラーが発生する
		fmt.Println(err)
		return
	}

	*/


	if err != nil {
		fmt.Println(err)
		return si.NewPostMessageInternalServerError().WithPayload(
			&si.PostMessageInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	defer file.Close()

	_, err = file.Write(p.Params.Image)
	if err != nil {
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
