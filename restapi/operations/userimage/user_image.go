package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"strconv"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

var assetsPath = os.Getenv("ASSETS_PATH")

func PostImage(p si.PostImagesParams) middleware.Responder {
	b64Img := p.Params.Image
	token := p.Params.Token

	tokenRepo := repositories.NewUserTokenRepository()
	imageRepo := repositories.NewUserImageRepository()

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return postImagesInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return postImagesUnauthorizedResponse()
	}

	id := tokenOwner.UserID

	// UPさせる画像をローカルに保存します。
	imagePath := assetsPath + "user" + strconv.Itoa(int(id))
	file, _ := os.Create(imagePath)
	defer file.Close()
	file.Write(b64Img)

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

/*			以下　Validationに用いる関数			*/

// 	プロフィール写真の更新
//	POST {hostname}/api/1.0/images
func postImagesInternalServerErrorResponse() middleware.Responder {
	return si.NewPostImagesInternalServerError().WithPayload(
		&si.PostImagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func postImagesUnauthorizedResponse() middleware.Responder {
	return si.NewPostImagesUnauthorized().WithPayload(
		&si.PostImagesUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

func postImagesBadRequestResponse() middleware.Responder {
	return si.NewPostImagesBadRequest().WithPayload(
		&si.PostImagesBadRequestBody{
			Code:    "401",
			Message: "Bad Request",
		})
}

func postImagesOKResponse(imagePath string) middleware.Responder {
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(imagePath),
		})
}
