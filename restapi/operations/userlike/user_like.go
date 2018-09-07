package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	tr := repositories.NewUserTokenRepository()
	lr := repositories.NewUserLikeRepository()
	ur := repositories.NewUserRepository()
	mr := repositories.NewUserMatchRepository()

	token, err := tr.GetByToken(p.Token)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error (at get token)",
			})
	}

	if token == nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code: "401",
				Message: "token is invalid",
			})
	}
	user, err := ur.GetByUserID(token.UserID)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error (at get token)",
			})
	}

	if user == nil {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	matchIds, err := mr.FindAllByUserID(user.ID)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error (at get token)",
			})
	}
	
	if p.Limit < 0 {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code: "400",
				Message: "limit is invalid",
			})
	}

	if p.Offset < 0 {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code: "400",
				Message: "offset is invalid",
			})
	}

	// マッチしていないlike
	var likes entities.UserLikes
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
		user, err :=  ur.GetByUserID(l.UserID)
		if err != nil {
			return si.NewGetLikesInternalServerError().WithPayload(
				&si.GetLikesInternalServerErrorBody{
					Code: "500",
					Message: "Internal Server Error",
				})
		}
		if user == nil {
			return si.NewGetLikesBadRequest().WithPayload(
				&si.GetLikesBadRequestBody{
					Code: "400",
					Message: "Bad Request",
				})
		}
		res.ApplyUser(*user)
		responses = append(responses, res)
	}

	sRes := responses.Build()
	return si.NewGetLikesOK().WithPayload(sRes)
}

// p.UserIDはpartnerのID
func PostLike(p si.PostLikeParams) middleware.Responder {
	tr := repositories.NewUserTokenRepository()
	ur := repositories.NewUserRepository()
	ulr := repositories.NewUserLikeRepository()
	mchr := repositories.NewUserMatchRepository()


	token, err := tr.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error(at get token)",
			})
	}

	if token == nil {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code: "401",
				Message: "invalid token",
			})
	}

	usr, err := ur.GetByUserID(token.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error(at get user)",
			})
	}

	if usr == nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	ptnr, err := ur.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error(at get user)",
			})
	}

	if ptnr == nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	like, err := ulr.GetLikeBySenderIDReceiverID(usr.ID, ptnr.ID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code: "500",
				Message: "Internal Server Error(at get user)",
			})
	}


	if like != nil {
		// 2回目
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code: "400",
				Message: "double like",
			})
	} else {

				
		if usr.ID == ptnr.ID {
			// 自分自身へのlikeなのでerror
			return si.NewPostLikeBadRequest().WithPayload(
				&si.PostLikeBadRequestBody{
					Code: "400",
					Message: "like yourself",
				})
		}
		
		if usr.Gender == ptnr.Gender {
			// 同性なのでerror
			return si.NewPostLikeBadRequest().WithPayload(
				&si.PostLikeBadRequestBody{
					Code: "400",
					Message: "same geneder",
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


		revlike, err := ulr.GetLikeBySenderIDReceiverID(ptnr.ID, usr.ID)
	
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
