package userimage

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"image"
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

// PUT Profile Image
func PostImage(p si.PostImagesParams) middleware.Responder {
	userimageHandler := repositories.NewUserImageRepository()
	usertokenHandler := repositories.NewUserTokenRepository()
	// find token
	usertokn, err := usertokenHandler.GetByToken(p.Params.Token)
	if err != nil {
		return RespInternalErr()
	}

	if usertokn == nil {
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// Define save directory
	writeimage := p.Params.Image
	convert := sha512.Sum512([]byte(usertokn.Token))
	pathdir := "assets/" + hex.EncodeToString(convert[:])
	checkext := bytes.NewReader(writeimage)
	_, extension, err := image.DecodeConfig(checkext)
	if err != nil {
		return RespInternalErr()
	}
	// check image extension
	switch extension {
	case "png":
		pathdir += ".png"
	case "jpg":
		pathdir += ".jpg"
	default:
		return RespInternalErr()
	}
	file, err := os.Create(pathdir)
	if err != nil {
		return RespInternalErr()
	}
	defer file.Close()
	//write picture
	file.Write(writeimage)
	userimage := entities.UserImage{
		UserID: usertokn.UserID,
		Path:   "https://si-2018-008.eure.jp/" + pathdir,
	}

	err = userimageHandler.Update(userimage)
	if err != nil {
		return RespInternalErr()
	}
	return PutImageOK("https://si-2018-008.eure.jp/" + pathdir)

}

// return 200 OK
func PutImageOK(savepath string) middleware.Responder {
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(savepath),
		})
}

// return 500 Internal Server Error
func RespInternalErr() middleware.Responder {
	return si.NewPostImagesInternalServerError().WithPayload(
		&si.PostImagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

// return 400 Bad Request
func RespBadReqestErr() middleware.Responder {
	return si.NewPostImagesBadRequest().WithPayload(
		&si.PostImagesBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}
