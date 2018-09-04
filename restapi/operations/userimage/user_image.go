package userimage

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"

	// "github.com/eure/si2018-server-side/entities"
	// "github.com/eure/si2018-server-side/repositories"
	
)

func PostImage(p si.PostImagesParams) middleware.Responder {
	// ur := repositories.NewUserRepository()
	// uir := repositories.NewUserImageRepository()

	// usr, _ := ur.GetByToken(p.Params.Token)
	

	return si.NewPostImagesOK()
}
