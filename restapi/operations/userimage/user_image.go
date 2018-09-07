package userimage

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"strconv"
	"bytes"
	"os"
	
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	i := repositories.NewUserImageRepository()
	t := repositories.NewUserTokenRepository()
	
	// Params から, Imageの取得
	uploadImage := p.Params.Image
	
	// tokenから UserToken entitiesの取得 (Validation)
	token := p.Params.Token
	loginUserToken,err := t.GetByToken(token)
	if err != nil {
		outPutPostStatus(500)
	}
	if loginUserToken == nil {
		outPutPostStatus(401)
	}
	
	// プロフィール画像の保存先情報の取得
	// local
	var localpath string
	localpath = os.Getenv("ASSETS_PATH")
	// db
	var dbpath string
	dbpath = os.Getenv("ASSETS_BASE_URI")
	
	// プロフィール画像の保存する際の名前
	fileName := strconv.Itoa(int(loginUserToken.UserID))
	
	// 取得したImageの形式を調べる
	leader := bytes.NewBuffer(uploadImage)
	_, format , err := image.DecodeConfig(leader)
	
	if err != nil {
		return outPutPostStatus(500)
	}
	if format != "jpeg" && format != "png" && format != "gif" {
		return outPutPostStatus(400)
	}
	
	// プロフィール画像の名前
	// 拡張子によって名前を変更
	// localに保存するプロフィール画像のPATH
	var filePathLocal string
	filePathLocal = localpath + fileName + "." + format
	// dbに保存するプロフィール画像のPATH
	var filePathDB string
	filePathDB = dbpath + fileName + "." + format
	
	// localに空のimage fileを用意
	file , err := os.Create(filePathLocal)
	defer file.Close()
	file.Write(uploadImage)
	
	var updatedProfile entities.UserImage
	updatedProfile.UserID    = loginUserToken.UserID
	updatedProfile.Path      = filePathDB
	
	err = i.Update(updatedProfile)
	
	if err != nil {
		outPutPostStatus(500)
	}
	
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(filePathDB),
		})
}

func outPutPostStatus (num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	case 401:
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗)",
			})
	case 400:
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	return nil
}