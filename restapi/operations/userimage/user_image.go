package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	//"github.com/eure/si2018-server-side/restapi/operations/util"
	"github.com/go-openapi/runtime/middleware"
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	/*token := p.Params.Token
	img := p.Params.Image

	r := repositories.NewUserImageRepository()

	res := si.PostImagesOKBody{}*/
	return si.NewPostImagesOK()
}
