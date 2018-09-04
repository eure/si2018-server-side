package userlike

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	userTokenEnt, err := repositories.NewUserTokenRepository().GetByToken(p.Token)

	if err != nil {
		return getLikesInternalServerErrorResponse()
	}

	if userTokenEnt == nil {
		return getLikesUnauthorizedResponse()
	}

	// int64になっているのでcastする必要がある
	limit := int(p.Limit)
	offset := int(p.Offset)
	userID := userTokenEnt.UserID

	userLikeRepository := repositories.NewUserLikeRepository()
	likesEnt, err := userLikeRepository.FindGotLikeWithLimitOffset(userID, limit, offset, nil)

	var likeUserResponsesEnt entities.LikeUserResponses

	for _, likeUserEnt := range likesEnt {
		likeResponse := entities.LikeUserResponse{LikedAt: likeUserEnt.CreatedAt}

		//FIXME ループ内でクエリ発行は最低の行為のような気がする
		user, err := repositories.NewUserRepository().GetByUserID(likeUserEnt.UserID)
		if err != nil {
			return getLikesInternalServerErrorResponse()

		}
		likeResponse.ApplyUser(*user)

		likeUserResponsesEnt = append(likeUserResponsesEnt, likeResponse)
	}

	likeResponses := likeUserResponsesEnt.Build()

	return si.NewGetLikesOK().WithPayload(likeResponses)
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
