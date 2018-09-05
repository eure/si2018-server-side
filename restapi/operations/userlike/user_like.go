package userlike

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	rt := repositories.NewUserTokenRepository()
	r := repositories.NewUserRepository()
	rl := repositories.NewUserLikeRepository()
	rm := repositories.NewUserMatchRepository()
	t, err := rt.GetByToken(p.Token)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: GetByToken failed: " + err.Error(),
			})
	}
	if t == nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗): GetByToken failed",
			})
	}
	matched, err := rm.FindAllByUserID(t.UserID)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: FindAllByUserID failed: " + err.Error(),
			})
	}
	like, err := rl.FindGotLikeWithLimitOffset(t.UserID, int(p.Limit), int(p.Offset), matched)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: FindLikeAll failed: " + err.Error(),
			})
	}
	ids := make([]int64, 0)
	for _, l := range like {
		ids = append(ids, l.UserID)
	}
	users, err := r.FindByIDs(ids)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error: FindByIDs failed: " + err.Error(),
			})
	}
	sEnt := make([]*models.LikeUserResponse, 0)
	for i, l := range like {
		response := entities.LikeUserResponse{LikedAt: l.UpdatedAt}
		response.ApplyUser(users[i])
		swaggerLike := response.Build()
		sEnt = append(sEnt, &swaggerLike)
	}

	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	// rt := repositories.NewUserTokenRepository()
	// r := repositories.NewUserRepository()
	// rl := repositories.NewUserLikeRepository()
	return si.NewPostLikeOK()
}
