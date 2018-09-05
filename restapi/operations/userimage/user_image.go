package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/restapi/operations/util"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	"os"
	"fmt"
	"strconv"
	"encoding/hex"
	"strings"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	token := p.Params.Token
	img := p.Params.Image //Image strfmt.Base64 -> type Base64 []byte

	err := util.ValidateToken(token)
	if err != nil {
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	id, _ := util.GetIDByToken(token)

	// Get filetype
	
	head := hex.EncodeToString(img[:4])

	var filetype string

	if strings.HasPrefix(head, "ffd8") {
		filetype = ".jpg"
	} else if head == "89504e47" {
		filetype = ".png"
	} else if head == "47494638" {
		filetype = ".gif"
	} else {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// Create image file
	ap := os.Getenv("ASSETS_PATH")
	filename := "icon" + strconv.Itoa(int(id)) + filetype

	f, err := os.Create(ap + filename)
	if err != nil {
		fmt.Println(err)
		return si.NewPostImagesOK()
	}
	defer f.Close()
	f.Write(img)

	// Update DB
	ui := entities.UserImage{}
	ui.UserID = id
	ui.Path = ap + filename

	r := repositories.NewUserImageRepository()
	err = r.Update(ui) /* TODO check existance? */

	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(ap + filename),
		})
}
