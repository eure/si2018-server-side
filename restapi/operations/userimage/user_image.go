package userimage

import (
	"errors"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	"github.com/eure/si2018-server-side/repositories"

	"github.com/eure/si2018-server-side/entities"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func postImageThrowInternalServerError(fun string, err error) *si.PostImagesInternalServerError {
	return si.NewPostImagesInternalServerError().WithPayload(
		&si.PostImagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func postImageThrowUnauthorized(mes string) *si.PostImagesUnauthorized {
	return si.NewPostImagesUnauthorized().WithPayload(
		&si.PostImagesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func postImageThrowBadRequest(mes string) *si.PostImagesBadRequest {
	return si.NewPostImagesBadRequest().WithPayload(
		&si.PostImagesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func getFreeFileName(path string) (string, error) {
	// 100000 枚までしか写真を保存しない
	const len = 8
	const format = "img%05d.png"
	const limit = 100000
	for i := 0; i < limit; i++ {
		filename := fmt.Sprintf(format, i)
		if !exists(path + filename) {
			return filename, nil
		}
	}
	return "", errors.New("storage capacity exceeded")
}

func PostImage(p si.PostImagesParams) middleware.Responder {
	var err error
	imageRepo := repositories.NewUserImageRepository()
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Params.Token)
		if err != nil {
			return postImageThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return postImageThrowUnauthorized("GetByToken failed")
		}
		id = token.UserID
	}
	// 既にある画像のタイムスタンプを使いまわしたい
	image, err := imageRepo.GetByUserID(id)
	if err != nil {
		return postImageThrowInternalServerError("GetByUserID", err)
	}
	var localFile, serverFile string
	{
		path := os.Getenv("ASSETS_PATH")
		uri := os.Getenv("ASSETS_BASE_URI")
		filename, err := getFreeFileName(path)
		if err != nil {
			return postImageThrowInternalServerError("getFreeFileName", err)
		}
		localFile = path + filename
		serverFile = uri + filename
	}
	// 画像をローカルに保存する
	{
		file, err := os.Create(localFile)
		if err != nil {
			return postImageThrowInternalServerError("Create", err)
		}
		// file.Close() でも error が吐かれる可能性がある
		// とりあえず今は対処しない
		defer file.Close()
		file.Write(p.Params.Image)
	}
	// DB を更新
	now := strfmt.DateTime(time.Now())
	ent := entities.UserImage{
		UserID:    id,
		Path:      serverFile,
		UpdatedAt: now}
	// 既に DB にデータがあれば Update、なければ Create
	if image == nil {
		err = imageRepo.Create(ent)
		if err != nil {
			return postImageThrowInternalServerError("Create", err)
		}
	} else {
		err = imageRepo.Update(ent)
		if err != nil {
			return postImageThrowInternalServerError("Update", err)
		}
	}
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(ent.Path)})
}
