package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	return si.NewPostImagesOK()
}
