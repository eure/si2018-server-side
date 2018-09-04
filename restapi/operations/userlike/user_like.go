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
	// マッチしていないlike
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

// p.UserIDはpartnerのID
func PostLike(p si.PostLikeParams) middleware.Responder {
	ur := repositories.NewUserRepository()
	ulr := repositories.NewUserLikeRepository()
	mchr := repositories.NewUserMatchRepository()

	usr, _ := ur.GetByToken(p.Params.Token)
	ptnr, _ := ur.GetByUserID(p.UserID)

	like, _ := ulr.GetLikeBySenderIDReceiverID(usr.ID, ptnr.ID)

	if like != nil {
		// 2回目
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code: "400",
				Message: "二回",
			})
	} else {
		if usr.Gender == ptnr.Gender {
			// 同性なのでerror
			return si.NewPostLikeBadRequest().WithPayload(
				&si.PostLikeBadRequestBody{
					Code: "400",
					Message: "同性",
				})
		}
		
		if usr.ID == ptnr.ID {
			// 自分自身へのlikeなのでerror
			return si.NewPostLikeBadRequest().WithPayload(
				&si.PostLikeBadRequestBody{
					Code: "400",
					Message: "ナルシスト",
				})
		}

		newLike := entities.NewUserLike(usr.ID, ptnr.ID)
		err := ulr.Create(newLike)
		if err != nil {
			return si.NewPostLikeInternalServerError().WithPayload(
				&si.PostLikeInternalServerErrorBody{
					Code: "500",
					Message: "Internal Server Error",
				})
		}


		revlike, _ := ulr.GetLikeBySenderIDReceiverID(ptnr.ID, usr.ID)
	
		if revlike != nil {
			// このタイミングでマッチしたことになるのでlikeCreateとmatchCreate
			newMatch := entities.NewUserMatch(usr.ID, ptnr.ID)
			err := mchr.Create(newMatch)
			if err != nil {
				return si.NewPostLikeInternalServerError().WithPayload(
					&si.PostLikeInternalServerErrorBody{
						Code: "500",
						Message: "Internal Server Error",
					})
			}
		}
	}


	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code: "200",
			Message: "OK",
		})


}
