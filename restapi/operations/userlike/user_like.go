package userlike

import (
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	UserMatchRepository := repositories.NewUserMatchRepository()
	UserLikeRepository := repositories.NewUserLikeRepository()
	UserRepository := repositories.NewUserRepository()
	UserTokenRepository := repositories.NewUserTokenRepository()
	UserByToken, err := UserTokenRepository.GetByToken(p.Token)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if UserByToken == nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	if p.Limit <= 0 {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	UserID := UserByToken.UserID
	ids, err := UserMatchRepository.FindAllByUserID(UserID)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	Likes, err := UserLikeRepository.FindGotLikeWithLimitOffset(UserID, int(p.Limit), int(p.Offset), ids)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	// いいねされていないとBadになってしまう
	/*
		if Likes == nil {
			return si.NewGetLikesBadRequest().WithPayload(
				&si.GetLikesBadRequestBody{
					Code:    "400",
					Message: "Bad Request",
				})
		}
	*/
	var LikeUserIDs []int64
	for _, Like := range Likes {
		LikeUserIDs = append(LikeUserIDs, Like.UserID)
	}
	LikeUsers, err := UserRepository.FindByIDs(LikeUserIDs)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	var LikeUserResponses entities.LikeUserResponses
	for _, like := range Likes {
		LikeUserResponse := entities.LikeUserResponse{}
		LikeUserResponse.LikedAt = like.CreatedAt
		for _, likeUser := range LikeUsers {
			if likeUser.ID == like.UserID {
				LikeUserResponse.ApplyUser(likeUser)
			}
		}
		LikeUserResponses = append(LikeUserResponses, LikeUserResponse)
	}
	sEnt := LikeUserResponses.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	UserTokenRepository := repositories.NewUserTokenRepository()
	UserLikeRepository := repositories.NewUserLikeRepository()
	UserMatchRepository := repositories.NewUserMatchRepository()
	Token := p.Params.Token
	UserByToken, err := UserTokenRepository.GetByToken(Token)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if UserByToken == nil {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}
	UserID := UserByToken.UserID
	PartnerID := p.UserID
	if PartnerID <= 0 {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	if UserID == PartnerID {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	//いいねする相手をいいねしていないか
	Like, err := UserLikeRepository.GetLikeBySenderIDReceiverID(UserID, PartnerID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if Like != nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	InsertLike := entities.UserLike{
		UserID:    UserID,
		PartnerID: PartnerID,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}
	err = UserLikeRepository.Create(InsertLike)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	//いいねした相手にいいねされているかどうか
	Like, err = UserLikeRepository.GetLikeBySenderIDReceiverID(PartnerID, UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if Like != nil {
		InsertMatch := entities.UserMatch{
			UserID:    PartnerID,
			PartnerID: UserID,
			CreatedAt: strfmt.DateTime(time.Now()),
			UpdatedAt: strfmt.DateTime(time.Now()),
		}
		err = UserMatchRepository.Create(InsertMatch)
	}
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:    "200",
			Message: "OK",
		})
}
