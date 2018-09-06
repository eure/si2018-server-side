package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"os"
	"time"
	"strconv"
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/strfmt"
)

func PostImage(p si.PostImagesParams) middleware.Responder {

	// レポジトリを初期化する
	tokenR := repositories.NewUserTokenRepository()
	userImageR := repositories.NewUserImageRepository()

	// 400エラー
	if p.Params.Token == "" {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message:  "Can't find token.",
			})
	}

	// トークンを検索する
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)

	// 401エラー
	if tokenEnt == nil {
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code: "401",
				Message: "Your token is invalid.",
			})
	}

	// 500エラー
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code: "500",
				Message: "Interval Server Error.",
			})
	}

	// ファイルデータから拡張子を判別する
	extension := findExtensionFromByteFile(p.Params.Image)

	if extension == "" {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code: "400",
				Message: "Bad Request. Only jpeg and png file can be accepted.",
			})
	}

	// 挿入するファイル名を決定する
	fileName := strconv.FormatInt(time.Now().Unix(), 10) + "." + extension

	// 挿入先のパスを指定する
	insertPath := os.Getenv("ASSETS_PATH") + fileName

	// ファイルに書き込みをする
	file, _ := os.Create(insertPath)

	defer file.Close()

	file.Write(p.Params.Image)

	imageUrl := "http://localhost:8080/assets/" + fileName

	// 挿入するための構造体を作成する
	userImageEnt := entities.UserImage{
		UserID: tokenEnt.UserID,
		Path: imageUrl,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}

	// 画像の情報を挿入する
	err = userImageR.Update(userImageEnt)

	// 500エラー
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}

	// 結果を返す
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(imageUrl),
		})
}

// 拡張子を判別する
func findExtensionFromByteFile(fileData []byte) string {
	fileDataBytes := []byte(fileData)

	var jpegBytes = []byte{'\xff', '\xd8'}
	var pngBytes = []byte{'\x89', '\x50', '\x4e', '\x47', '\x0d', '\x0a', '\x1a', '\x0a'}

	if fileDataBytes[0] == jpegBytes[0] && fileDataBytes[1] == jpegBytes[1] {
		return "jpeg"
	}

	if fileDataBytes[0] == pngBytes[0] &&
		fileDataBytes[1] == pngBytes[1] &&
		fileDataBytes[2] == pngBytes[2] &&
		fileDataBytes[3] == pngBytes[3] &&
		fileDataBytes[4] == pngBytes[4] &&
		fileDataBytes[5] == pngBytes[5] &&
		fileDataBytes[6] == pngBytes[6] &&
		fileDataBytes[7] == pngBytes[7] {
			return "png"
	}

	return ""
}