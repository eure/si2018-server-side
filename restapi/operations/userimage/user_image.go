package userimage

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/nfnt/resize"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	i := repositories.NewUserImageRepository()

	// ユーザーID取得用
	token, _ := t.GetByToken(p.Params.Token)

	var assetsPath string
	var assetsBaseURI string

	// assetsのenvを取得
	assetsPath = os.Getenv("ASSETS_PATH")
	assetsBaseURI = os.Getenv("ASSETS_BASE_URI")

	// 変数に代入
	img := p.Params.Image

	// 画像のサイズを調べる(1M以上のときはBad Request)
	// 画像の容量が大きすぎるものをあげさせないようにするため
	reader := bytes.NewBuffer(img)
	size := int64(len(img))
	if size >= 1000000 {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// 画像のリサイズ
	// 全て300x300の正方形に統一
	// フォーマットも確認しておく
	decode, format, err := image.Decode(reader)
	if err != nil {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	// resize.Lanczos3は圧縮アルゴリズム
	resized := resize.Resize(300, 300, decode, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = imageEncode(buf, resized, format)
	if err != nil {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// pathの名前を定義
	pathname := assetsBaseURI + p.Params.Token + "." + format
	// fileの名前を定義
	filename := assetsPath + p.Params.Token + "." + format

	// ファイル作成
	f, err := os.Create(filename)
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	defer f.Close()

	// ファイルに書き込む
	f.Write(buf.Bytes())

	// アップデートしたい値の定義
	init := entities.UserImage{
		Path: pathname,
	}
	// 更新
	err = i.Update(init)
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

func imageEncode(target io.Writer, image image.Image, format string) error {
	switch format {
	case "jpeg", "jpg":
		jpeg.Encode(target, image, nil)
	case "png":
		png.Encode(target, image)
	case "gif":
		gif.Encode(target, image, nil)
	default:
		return errors.New("invalid format")
	}
	return nil
}
