package userlike

import (
	"fmt"

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
	userTokenEnt, err := repositories.NewUserTokenRepository().GetByToken(p.Params.Token)

	if err != nil {
		return getLikesInternalServerErrorResponse()
	}

	if userTokenEnt == nil {
		return getLikesUnauthorizedResponse()
	}

	userLike := entities.UserLike{
		UserID:    userTokenEnt.UserID,
		PartnerID: p.UserID,
	}

	userLikeRepository := repositories.NewUserLikeRepository()
	errs := userLikeRepository.Validate(userLike)
	if errs != nil {
		str := fmt.Sprintf("%v", errs)
		return postLikeBadRequestResponse(str)
	}

	err = userLikeRepository.Create(userLike)

	if err != nil {
		return postLikeInternalServerErrorResponse()
	}

	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:    "200",
			Message: "OK",
		})
}
