package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	var assetsPath string
	var assetsBaseURI string

	// assetsのenvを取得
	assetsPath = os.Getenv("ASSETS_PATH")
	assetsBaseURI = os.Getenv("ASSETS_BASE_URI")

	// 変数に代入
	img := p.Params.Image

	t := repositories.NewUserTokenRepository()
	i := repositories.NewUserImageRepository()

	// ユーザーID取得用
	token, _ := t.GetByToken(p.Params.Token)

	// pathの名前を定義
	pathname := assetsBaseURI + p.Params.Token + ".png"
	// fileの名前を定義
	filename := assetsPath + p.Params.Token + ".png"

	// エラー処理を書くべきか調べる
	// ファイル作成
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(img)

	// アップデートしたい値の定義
	init := entities.UserImage{
		Path: pathname,
	}
	// 更新
	err := i.Update(init)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	// 更新した値からpathを取得するために
	ent, err := i.GetByUserID(token.UserID)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if ent == nil {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(ent.Path),
		})
}
