package userimage

import (
	"bytes"
	"fmt"
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
	nuir := repositories.NewUserImageRepository()
	nutr := repositories.NewUserTokenRepository()

	usertkn, err := nutr.GetByToken(p.Params.Token)
	if err != nil {
		fmt.Println("token err")
	}

	// Define save directory
	writeimage := p.Params.Image

	pathdir := "assets/" + usertkn.Token
	checkext := bytes.NewReader(writeimage)
	_, extension, err := image.DecodeConfig(checkext)
	if err != nil {
		fmt.Println(err)
	}
	// check image extension
	switch extension {
	case "png":
		pathdir += ".png"
	case "jpg":
		pathdir += ".jpg"
	}
	file, err := os.Create(pathdir)
	if err != nil {
		fmt.Println(err)
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

	err = nuir.Update(userimage)
	if err != nil {
		fmt.Println(err)
	}
	userimg, err := nuir.GetByUserID(usertkn.UserID)
	if err != nil {
		fmt.Println("not find your account")
	}
	fmt.Println(userimg)
	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(pathdir),
		})
}
