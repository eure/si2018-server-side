package userlike

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	return si.NewGetLikesOK()
	userTokenEnt, err := repositories.NewUserTokenRepository().GetByToken(p.Token)

	if err != nil {
		return getLikesInternalServerErrorResponse()
	}

	if userTokenEnt == nil {
		return getLikesUnauthorizedResponse()
	}

}

func PostLike(p si.PostLikeParams) middleware.Responder {
	return si.NewPostLikeOK()
}

func getLikesInternalServerErrorResponse() middleware.Responder {
	return si.NewGetLikesInternalServerError().WithPayload(
		&si.GetLikesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getLikesUnauthorizedResponse() middleware.Responder {
	return si.NewGetLikesUnauthorized().WithPayload(
		&si.GetLikesUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}
