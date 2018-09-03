package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	lr := repositories.NewUserLikeRepository()
	ur := repositories.NewUserRepository()
	mr := repositories.NewUserMatchRepository()
	user, _ := ur.GetByToken(p.Token)
	matchIds, _ := mr.FindAllByUserID(user.ID)
	
	var likes entities.UserLikes
	var err error
	likes, err = lr.FindGotLikeWithLimitOffset(user.ID, int(p.Limit), int(p.Offset), matchIds)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error",
			})
	}
	var responses entities.LikeUserResponses
	for _, l := range likes {
		var res = entities.LikeUserResponse{}
		user, _ :=  ur.GetByUserID(l.UserID)
		res.ApplyUser(*user)
		responses = append(responses, res)
	}

	sRes := responses.Build()
	return si.NewGetLikesOK().WithPayload(sRes)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	return si.NewPostLikeOK()
}
