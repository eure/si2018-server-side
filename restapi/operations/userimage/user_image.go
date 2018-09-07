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
	"log"
	//"strconv"
	"encoding/hex"
	"strings"
	"math/rand"
	"time"
	"path/filepath"
)

const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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

	rand.Seed(time.Now().UnixNano())
	const filenamelen = 8

	filename := getRandStr(filenamelen) + filetype

	f, err := os.Create(filepath.Join(ap, filename)) // Permission 0666
	if err != nil {
		log.Print("File create failed")
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	defer f.Close()

	_, err = f.Write(img)
	if err != nil {
		log.Print("File write failed")
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// Update DB
	ui := entities.UserImage{}
	ui.UserID = id
	ui.Path = ap + filename

	r := repositories.NewUserImageRepository()

	old, _ := r.GetByUserID(id) /* TODO check existance? */
	log.Print(old)

	err = r.Update(ui) 
	if err != nil {
		return si.NewPostImagesInternalServerError().WithPayload(
			&si.PostImagesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	new, _ := r.GetByUserID(id)
	log.Print(new)

	return si.NewPostImagesOK().WithPayload(
		&si.PostImagesOKBody{
			ImageURI: strfmt.URI(filepath.Join(ap, filename)),
		})
}

func getRandStr(n int) string {
    var s []byte
    for i := 0; i < n; i++ {
        s = append(s, letters[rand.Intn(len(letters))])
    }
    return string(s)
}
