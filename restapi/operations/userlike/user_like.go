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
	var ids []int64

	for _, likeEnt := range likesEnt {
		ids = append(ids, likeEnt.UserID)

		likeResponse := entities.LikeUserResponse{LikedAt: likeEnt.CreatedAt}
		likeUserResponsesEnt = append(likeUserResponsesEnt, likeResponse)
	}

	var users entities.Users
	users, err = repositories.NewUserRepository().FindByIDs(ids)
	if err != nil {
		return getLikesInternalServerErrorResponse()
	}

	likeUserResponsesEnt = likeUserResponsesEnt.ApplyUsers(users)
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

func postLikeInternalServerErrorResponse() middleware.Responder {
	return si.NewPostLikeInternalServerError().WithPayload(
		&si.PostLikeInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func postLikeUnauthorizedResponse() middleware.Responder {
	return si.NewPostLikeUnauthorized().WithPayload(
		&si.PostLikeUnauthorizedBody{
			Code:    "401",
			Message: "Token Is Invalid",
		})
}

func postLikeBadRequestResponse(message string) middleware.Responder {
	return si.NewPostLikeBadRequest().WithPayload(
		&si.PostLikeBadRequestBody{
			Code:    "400",
			Message: message,
		})
}
