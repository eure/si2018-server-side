package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"fmt"
	"os"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	tokenR := repositories.NewUserTokenRepository()
	tokenEnt, err := tokenR.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}

	// 元のimageのパスを取得する -> 必要ない気が..
	imageR := repositories.NewUserImageRepository()
	imageEnt, err := imageR.GetByUserID(tokenEnt.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code    : "500",
				Message : "Internal Server Error",
			})
	}
	fmt.Println(imageEnt.Path)

	file, _ := os.Open(imageEnt.Path)
	fmt.Println(file)


	// Base64データから画像ファイルを作成して、
	// ランダムな名前で保存する -> ユーザIDを含めたPATHにして初回以降はこれを更新する形にする？
	// パスの情報を更新する Update
	// 正常に終了したら

	return si.NewPostImagesOK()
}
