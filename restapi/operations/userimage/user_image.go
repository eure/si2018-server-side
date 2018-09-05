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
	// 100000000 枚までしか写真を保存しない
	const len = 8
	const format = "img%08d.png"
	var filename string
	for i := 0; i < 100000000; i++ {
		filename = fmt.Sprintf(format, i)
		if !exists(path + filename) {
			break
		}
	}
	if filename == "" {
		return "", errors.New("storage capacity exceeded")
	}
	return filename, nil
}

func PostImage(p si.PostImagesParams) middleware.Responder {
	ri := repositories.NewUserImageRepository()
	rt := repositories.NewUserTokenRepository()
	t, err := rt.GetByToken(p.Params.Token)
	if err != nil {
		return postImageThrowInternalServerError("GetByToken", err)
	}
	if t == nil {
		return postImageThrowUnauthorized("GetByToken failed")
	}
	image, err := ri.GetByUserID(t.UserID)
	if err != nil {
		return postImageThrowInternalServerError("GetByUserID", err)
	}
	path := os.Getenv("ASSETS_PATH")
	uri := os.Getenv("ASSETS_BASE_URI")
	filename, err := getFreeFileName(path)
	if err != nil {
		return postImageThrowInternalServerError("getFreeFileName", err)
	}
	file, err := os.Create(path + filename)
	if err != nil {
		return postImageThrowInternalServerError("Create", err)
	}
	// file.Close() でも error が吐かれる可能性がある
	// とりあえず今は対処しない
	defer file.Close()
	file.Write(p.Params.Image)
	now := strfmt.DateTime(time.Now())
	fmt.Printf("Hello! path: %v, id: %d", path, t.UserID)
	ent := entities.UserImage{
		UserID:    t.UserID,
		Path:      uri + filename,
		UpdatedAt: now}
	if image == nil {
		err = ri.Create(ent)
		if err != nil {
			return postImageThrowInternalServerError("Create", err)
		}
	} else {
		err = ri.Update(ent)
		if err != nil {
			return postImageThrowInternalServerError("Update", err)
		}
	}
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(ent.Path)})
}
