package userimage

import (
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
	// idからUserImageを持ってくる
	// バイナリーで送られてきた写真をローカルに保存
	// ファイル名も持ってこれるの?
	// UserImggeに写真のパスを入れる
	// バイナリからファイルの形式をとる
	repUserImage := repositories.NewUserImageRepository()
	repUserToken := repositories.NewUserTokenRepository()

	loginUser, _ := repUserToken.GetByToken(p.Params.Token)

	assetsPath := os.Getenv("ASSETS_PATH")
	imagePath := assetsPath + "user" + strconv.Itoa(int(loginUser.UserID))

	file, _ := os.Create(imagePath)
	defer file.Close()

	file.Write(p.Params.Image)

	userImage, _ := repUserImage.GetByUserID(loginUser.UserID)

	userImage.Path = imagePath
	userImage.UpdatedAt = strfmt.DateTime(time.Now())

	err := repUserImage.Update(*userImage)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	updatedUserImage, _ := repUserImage.GetByUserID(loginUser.UserID)

	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(updatedUserImage.Path),
		})
}
