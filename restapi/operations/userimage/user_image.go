package userimage

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	userimageHandler := repositories.NewUserImageRepository()
	usertokenHandler := repositories.NewUserTokenRepository()

	usertkn, err := usertokenHandler.GetByToken(p.Params.Token)
	if err != nil {
		return RespInternalErr()
	}

	if usertkn == nil {
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// Define save directory
	writeimage := p.Params.Image

	pathdir := "assets/" + usertkn.Token
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
		UserID:    usertkn.UserID,
		Path:      pathdir,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}

	err = userimageHandler.Update(userimage)
	if err != nil {
		return RespInternalErr()
	}
	return PutImageOK(pathdir)

}

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

func RespBadReqestErr() middleware.Responder {
	return si.NewPostImagesBadRequest().WithPayload(
		&si.PostImagesBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}
