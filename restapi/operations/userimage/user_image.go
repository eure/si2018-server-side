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

	// 挿入するファイル名を決定する
	fileName := strconv.FormatInt(time.Now().Unix(), 10) + ".png"

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
