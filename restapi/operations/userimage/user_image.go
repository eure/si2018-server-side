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
	token := p.Params.Token

	t := repositories.NewUserTokenRepository()
	i := repositories.NewUserImageRepository()

	// ユーザーID取得用
	ut, _ := t.GetByToken(token)
	// pathの名前を定義
	pathname := assetsBaseURI + token + ".png"
	// fileの名前を定義
	filename := assetsPath + token + ".png"
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
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	// 更新した値からpathを取得するために
	ent, _ := i.GetByUserID(ut.UserID)

	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(ent.Path),
		})
}
