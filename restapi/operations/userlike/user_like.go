package userlike

import (
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func getLikesThrowInternalServerError(fun string, err error) *si.GetLikesInternalServerError {
	return si.NewGetLikesInternalServerError().WithPayload(
		&si.GetLikesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func getLikesThrowUnauthorized(mes string) *si.GetLikesUnauthorized {
	return si.NewGetLikesUnauthorized().WithPayload(
		&si.GetLikesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func getLikesThrowBadRequest(mes string) *si.GetLikesBadRequest {
	return si.NewGetLikesBadRequest().WithPayload(
		&si.GetLikesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetLikes(p si.GetLikesParams) middleware.Responder {
	rt := repositories.NewUserTokenRepository()
	r := repositories.NewUserRepository()
	rl := repositories.NewUserLikeRepository()
	rm := repositories.NewUserMatchRepository()
	t, err := rt.GetByToken(p.Token)
	if err != nil {
		return getLikesThrowInternalServerError("GetByToken", err)
	}
	if t == nil {
		return getLikesThrowUnauthorized("GetByToken failed")
	}
	matched, err := rm.FindAllByUserID(t.UserID)
	if err != nil {
		return getLikesThrowInternalServerError("FindAllByUserID", err)
	}
	like, err := rl.FindGotLikeWithLimitOffset(t.UserID, int(p.Limit), int(p.Offset), matched)
	if err != nil {
		return getLikesThrowInternalServerError("FindGotLikeWithLimitOffset", err)
	}
	ids := make([]int64, 0)
	for _, l := range like {
		ids = append(ids, l.UserID)
	}
	users, err := r.FindByIDs(ids)
	if err != nil {
		return getLikesThrowInternalServerError("FindByIDs", err)
	}
	if len(users) != len(like) {
		return getLikesThrowBadRequest("FindByIDs failed")
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

func postLikeThrowInternalServerError(fun string, err error) *si.PostLikeInternalServerError {
	return si.NewPostLikeInternalServerError().WithPayload(
		&si.PostLikeInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func postLikeThrowUnauthorized(mes string) *si.PostLikeUnauthorized {
	return si.NewPostLikeUnauthorized().WithPayload(
		&si.PostLikeUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func postLikeThrowBadRequest(mes string) *si.PostLikeBadRequest {
	return si.NewPostLikeBadRequest().WithPayload(
		&si.PostLikeBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	rt := repositories.NewUserTokenRepository()
	r := repositories.NewUserRepository()
	rl := repositories.NewUserLikeRepository()
	rm := repositories.NewUserMatchRepository()
	t, err := rt.GetByToken(p.Params.Token)
	if err != nil {
		return postLikeThrowInternalServerError("GetByToken", err)
	}
	if t == nil {
		return postLikeThrowUnauthorized("GetByToken failed")
	}
	user, err := r.GetByUserID(t.UserID)
	if err != nil {
		return postLikeThrowInternalServerError("GetByUserID", err)
	}
	if user == nil {
		return postLikeThrowBadRequest("GetByUserID failed")
	}
	partner, err := r.GetByUserID(p.UserID)
	if err != nil {
		return postLikeThrowInternalServerError("GetByUserID", err)
	}
	if partner == nil {
		return postLikeThrowBadRequest("GetByUserID failed")
	}
	if user.Gender != partner.GetOppositeGender() {
		return postLikeThrowBadRequest("genders must be opposite")
	}
	like, err := rl.GetLikeBySenderIDReceiverID(user.ID, partner.ID)
	if err != nil {
		return postLikeThrowInternalServerError("GetLikeBySenderIDReceiverID", err)
	}
	if like != nil {
		return postLikeThrowBadRequest("like action duplicates")
	}
	reverse, err := rl.GetLikeBySenderIDReceiverID(partner.ID, user.ID)
	if err != nil {
		return postLikeThrowInternalServerError("GetLikeBySenderIDReceiverID", err)
	}
	now := strfmt.DateTime(time.Now())
	*like = entities.UserLike{
		UserID:    user.ID,
		PartnerID: partner.ID,
		CreatedAt: now,
		UpdatedAt: now}
	// like を書き込んだあと, match を書き込むときにエラーが発生すると致命的
	// https://qiita.com/komattio/items/838ea5df68eb076e8099
	// transaction を利用してまとめて書きこむ必要がある
	err = rl.Create(*like)
	if err != nil {
		return postLikeThrowInternalServerError("Create", err)
	}
	if reverse != nil {
		match := entities.UserMatch{
			UserID:    partner.ID,
			PartnerID: user.ID,
			CreatedAt: now,
			UpdatedAt: now}
		err = rm.Create(match)
		if err != nil {
			return postLikeThrowInternalServerError("Create", err)
		}
	}
	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:    "200",
			Message: "OK",
		})
}
