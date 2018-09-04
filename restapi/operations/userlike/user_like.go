package userlike

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	user_m_r := repositories.NewUserMatchRepository()
	user_l_r := repositories.NewUserLikeRepository()
	user_r := repositories.NewUserRepository()
	user_t_r := repositories.NewUserTokenRepository()
	userByToken, err := user_t_r.GetByToken(p.Token)
	// メンターにバリデーション聞く
	// 値が足りてない
	UserID := userByToken.UserID
	ids, err := user_m_r.FindAllByUserID(UserID)
	if err != nil {
		return si.NewGetMatchesInternalServerError().WithPayload(
			&si.GetMatchesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	ent, err := user_l_r.FindGotLikeWithLimitOffset(UserID, int(p.Limit), int(p.Offset), ids)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	var LikeUserIDs []int64
	for _, value := range ent {
		LikeUserIDs = append(ids, value.UserID)
	}
	LikeUsers, _ := user_r.FindByIDs(LikeUserIDs)
	// UserLikeからLikeUserResponsesへのコンバート
	var LikeUserResponses entities.LikeUserResponses
	for _, value := range LikeUsers {
		LikeUserResponse := entities.LikeUserResponse{}
		LikeUserResponse.ApplyUser(value)
		LikeUserResponses = append(LikeUserResponses, LikeUserResponse)
	}
	sEnt := LikeUserResponses.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	return si.NewPostLikeOK()
}
