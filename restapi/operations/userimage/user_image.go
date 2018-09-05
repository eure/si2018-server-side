package userimage

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/strfmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"time"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

//- プロフィール写真の更新
//- POST {hostname}/api/1.0/images
//- TokenのValidation処理を実装してください
//- プロフィール写真を更新してください
func PostImage(p si.PostImagesParams) middleware.Responder {
	repUserImage := repositories.NewUserImageRepository()
	repUserToken := repositories.NewUserTokenRepository()

	token := p.Params.Token

	// tokenのバリデーション
	err := repUserToken.ValidateToken(token)
	if err != nil {
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code: "401",
				Message: "Token Is Invalid",
			})
	}

	// tokenからuserTokenを取得
	loginUser, _ := repUserToken.GetByToken(token)

	// 画像ファイルのパス設定
	assetsPath := os.Getenv("ASSETS_PATH")
	imagePath := assetsPath + "user" + strconv.Itoa(int(loginUser.UserID))

	// パスからファイルを作成
	file, err := os.Create(imagePath)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	defer file.Close()

	// 作成したファイルに画像ファイルを書き込む
	file.Write(p.Params.Image)

	// 更新用のUserImageを作成
	var userImage entities.UserImage
	userImage.Path = imagePath
	userImage.UpdatedAt = strfmt.DateTime(time.Now()) //型を合わせる

	// UserImageを更新
	err = repUserImage.Update(userImage)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 更新したUserImageを取得
	updatedUserImage, err := repUserImage.GetByUserID(loginUser.UserID)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 更新したUserImageのpathを返す
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(updatedUserImage.Path),
		})
}
