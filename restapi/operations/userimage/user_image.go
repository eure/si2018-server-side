package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"encoding/hex"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

var (
	assetsPath = os.Getenv("ASSETS_PATH")
	tokenRepo  = repositories.NewUserTokenRepository()
	imageRepo  = repositories.NewUserImageRepository()
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	b64Img := p.Params.Image
	token := p.Params.Token

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return postImagesInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return postImagesUnauthorizedResponse()
	}

	id := tokenOwner.UserID

	// 画像ファイルのフォーマットをバイナリデータ処理で判別
	header := hex.EncodeToString(b64Img[:4])
	// 画像ファイルの先頭には、そのファイルがどのような種類の画像フォーマットであるかを指し示すデータが含まれています。
	// この先頭部分をチェックすることによって、画像フォーマットを判別することができます。
	var imageType string

	switch header {
	case "ffd8":
		imageType = ".jpg"
	case "89504e47":
		imageType = ".png"
	case "47494638":
		imageType = ".gif"
	}

	// UPさせる画像をローカルに保存します。
	imagePath := assetsPath + "User" + strconv.Itoa(int(id)) + imageType
	file, _ := os.Create(imagePath)
	defer file.Close()
	_, err = file.Write(b64Img)

	// ユーザーのプロフィール画像を取得します。
	userImage, err := imageRepo.GetByUserID(id)
	if err != nil {
		return postImagesInternalServerErrorResponse()
	}

	userImage.Path = imagePath

	err = imageRepo.Update(*userImage)
	if err != nil {
		return postImagesInternalServerErrorResponse()
	}

	return postImagesOKResponse(imagePath)
}
