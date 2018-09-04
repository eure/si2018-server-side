package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	"fmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/strfmt"
)

// いいね！表示API
func GetLikes(p si.GetLikesParams) middleware.Responder {
	tokenRepo := repositories.NewUserTokenRepository()
	likeRepo := repositories.NewUserLikeRepository()
	matchRepo := repositories.NewUserMatchRepository()
	userRepo := repositories.NewUserRepository()

	// tokenが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(p.Token)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if tokenOwner == nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// トークンの持ち主とマッチングしているお相手を取得します
	matches, err := matchRepo.FindAllByUserID(tokenOwner.UserID)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// トークンの持ち主とマッチ済みのお相手を除き、いいね！を送ってくれたお相手を取得します。
	Likes, err := likeRepo.FindGotLikeWithLimitOffset(tokenOwner.UserID, int(p.Limit), int(p.Offset), matches)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// いいね！を送ってくれたお相手のプロフィールを取得します。
	var ids []int64
	var likeAt []strfmt.DateTime
	for _, like := range Likes {
		likeAt = append(likeAt, like.CreatedAt)
		ids = append(ids, like.UserID)
	}
	fmt.Println(ids, likeAt)

	likeUsers, err := userRepo.FindByIDs(ids)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// models.LikeUserResponse と entities.User をマッピングします。
	var likeUserResponses []*models.LikeUserResponse
	for i, likeUser := range likeUsers {
		likeUserResponses = append(likeUserResponses, &models.LikeUserResponse{
			AnnualIncome:   likeUser.AnnualIncome,
			Birthday:       likeUser.Birthday,
			BodyBuild:      likeUser.BodyBuild,
			Child:          likeUser.Child,
			CostOfDate:     likeUser.CostOfDate,
			CreatedAt:      likeUser.CreatedAt,
			Drinking:       likeUser.Drinking,
			Education:      likeUser.Education,
			Gender:         likeUser.Gender,
			Height:         likeUser.Height,
			Holiday:        likeUser.Holiday,
			HomeState:      likeUser.HomeState,
			Housework:      likeUser.Housework,
			HowToMeet:      likeUser.HowToMeet,
			ID:             likeUser.ID,
			ImageURI:       likeUser.ImageURI,
			Introduction:   likeUser.Introduction,
			Job:            likeUser.Job,
			LikedAt:        likeAt[i],
			MaritalStatus:  likeUser.MaritalStatus,
			Nickname:       likeUser.Nickname,
			NthChild:       likeUser.NthChild,
			ResidenceState: likeUser.ResidenceState,
			Smoking:        likeUser.Smoking,
			Tweet:          likeUser.Tweet,
			UpdatedAt:      likeUser.UpdatedAt,
			WantChild:      likeUser.WantChild,
			WhenMarry:      likeUser.WhenMarry,
		})
		fmt.Println(likeUser.ID, likeAt[i])
	}

	return si.NewGetLikesOK().WithPayload(likeUserResponses)
}

// いいね！送信API
func PostLike(p si.PostLikeParams) middleware.Responder {
	tokenRepo := repositories.NewUserTokenRepository()
	likeRepo := repositories.NewUserLikeRepository()
	userRepo := repositories.NewUserRepository()
	matchRepo := repositories.NewUserMatchRepository()

	// tokenが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if tokenOwner == nil {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// ユーザーが異性かどうかを検証します。
	user, err := userRepo.GetByUserID(tokenOwner.UserID)
	partner, err := userRepo.GetByUserID(p.UserID)
	if user.Gender == partner.Gender {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// トークンの持ち主から指定したユーザーへいいね！を送信します
	like := entities.UserLike{
		UserID:    tokenOwner.UserID,
		PartnerID: p.UserID,
	}
	sendLike := likeRepo.Create(like)
	if sendLike != nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// いいね！を送信したお相手が自分にいいね！をしていたかどうかを検証します。
	liked, err := likeRepo.GetLikeBySenderIDReceiverID(like.PartnerID, like.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 過去にいいね！されていたらマッチングします。
	if liked != nil {
		mathe := entities.UserMatch{
			UserID:    liked.UserID,
			PartnerID: liked.PartnerID,
		}
		matching := matchRepo.Create(mathe)
		if matching != nil {
			return si.NewPostLikeInternalServerError().WithPayload(
				&si.PostLikeInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}

		return si.NewPostLikeOK().WithPayload(&si.PostLikeOKBody{
			Code:    "200",
			Message: "OK",
		})
	}

	return si.NewPostLikeOK().WithPayload(&si.PostLikeOKBody{
		Code:    "200",
		Message: "OK",
	})
}
