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

func getFreeFileName(path, extension string) (string, error) {
	now := time.Now().UnixNano()
	// 時刻がバレないように簡易的に暗号化する
	now = now ^ 123456789123456789
	// 同じ時刻では 100 枚までしか写真を保存しない
	// 1 ns 単位で重複することはほぼないはずなので, サーバーの時刻を止めない限り問題ない
	const limit = 100
	format := "img%x%02d" + extension
	for i := 0; i < limit; i++ {
		filename := fmt.Sprintf(format, now, i)
		if !exists(path + filename) {
			return filename, nil
		}
	}
	return "", errors.New("storage capacity exceeded")
}

// image のタグが tag と一致するか
func tagMatches(tag []byte, image strfmt.Base64) bool {
	seq := []byte(image)
	if len(seq) < len(tag) {
		return false
	}
	for i, t := range tag {
		if seq[i] != t {
			return false
		}
	}
	return true
}

var pngTag = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}

// image のファイル形式が png であるか
func isPngFile(image strfmt.Base64) bool {
	return tagMatches(pngTag, image)
}

var jpgTag = []byte{0xff, 0xd8}

// image のファイル形式が jpg であるか
func isJpgFile(image strfmt.Base64) bool {
	return tagMatches(jpgTag, image)
}

// image が png, jpg 形式であれば対応する拡張子 ".png" などを返す
func getFileType(image strfmt.Base64) (string, error) {
	if isPngFile(image) {
		return ".png", nil
	}
	if isJpgFile(image) {
		return ".jpg", nil
	}
	return "", errors.New("unknown file type")
}

// DB アクセス: 3 回
// 計算量: O(1)
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
		extension, err := getFileType(p.Params.Image)
		if err != nil {
			return postImageThrowInternalServerError("getFileType", err)
		}
		filename, err := getFreeFileName(path, extension)
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
