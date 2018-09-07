package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

var (
	tokenRepo = repositories.NewUserTokenRepository()
	likeRepo  = repositories.NewUserLikeRepository()
	matchRepo = repositories.NewUserMatchRepository()
	userRepo  = repositories.NewUserRepository()
	imageRepo = repositories.NewUserImageRepository()
)

// いいね！表示API
func GetLikes(p si.GetLikesParams) middleware.Responder {
	// 入力値のValidation処理をします。
	limit := int(p.Limit)
	if limit <= 0 {
		return getLikesBadRequestResponse()
	}

	offset := int(p.Offset)
	if offset < 0 {
		return getLikesBadRequestResponse()
	}

	token := p.Token

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return getLikesInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return getLikesUnauthorizedResponse()
	}

	id := tokenOwner.UserID

	// ユーザーとマッチングしているお相手を取得します。
	matches, err := matchRepo.FindAllByUserID(id)
	if err != nil {
		return getLikesInternalServerErrorResponse()
	}

	// ユーザーとマッチ済みのお相手を除き、いいね！を送ってくれたお相手を取得します。
	likes, err := likeRepo.FindGotLikeWithLimitOffset(id, limit, offset, matches)
	if err != nil {
		return getLikesInternalServerErrorResponse()
	}

	// ユーザーをいいねしているお相手のIDの配列を取得します。
	var ids []int64
	for _, like := range likes {
		ids = append(ids, like.UserID)
	}

	// idの配列からユーザーを取得します。
	users, err := userRepo.FindByIDs(ids)
	if err != nil {
		return getLikesInternalServerErrorResponse()
	}

	// ユーザーのプロフィール画像を取得します。
	images, err := imageRepo.GetByUserIDs(ids)
	if err != nil {
		return getLikesInternalServerErrorResponse()
	}

	// 取得したお相手のIDとユーザーをいいね！しているお相手のIDを比較し、マッピングします。
	var ents entities.LikeUserResponses
	for _, like := range likes {
		ent := entities.LikeUserResponse{}
		for _, image := range images {
			for _, user := range users {
				if like.UserID == user.ID && image.UserID == user.ID {
					ent.ApplyUser(user)
					ent.LikedAt = like.CreatedAt
					ent.ImageURI = image.Path
				}
			}
		}

		ents = append(ents, ent)
	}

	sEnt := ents.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

// いいね！送信API
func PostLike(p si.PostLikeParams) middleware.Responder {
	// 入力値のValidation処理をします。
	partnerID := p.UserID
	if partnerID <= 0 {
		return postLikeBadRequestResponses()
	}

	token := p.Params.Token

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return postLikeInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return postLikeUnauthorizedResponse()
	}

	id := tokenOwner.UserID

	// ユーザーが異性かどうかを検証します。
	user, err := userRepo.GetByUserID(id)
	if err != nil {
		return postLikeInternalServerErrorResponse()
	} else if user == nil {
		return postLikeBadRequestResponses()
	}

	partner, err := userRepo.GetByUserID(p.UserID)
	if err != nil {
		return postLikeInternalServerErrorResponse()
	} else if partner == nil {
		return postLikeBadRequestResponses()
	}

	// ユーザーとお相手が同性の場合は除きます。
	if user.Gender == partner.Gender {
		return postLikeBadRequestResponses()
	}

	// ユーザーが過去にお相手へ、いいね！をしていたかどうかを検証します。
	liked, err := likeRepo.GetLikeBySenderIDReceiverID(id, partnerID)
	if err != nil {
		return postLikeInternalServerErrorResponse()
	}
	// 過去にいいねしていた場合、除きます。
	if liked != nil {
		return postLikeBadRequestResponses()
	}

	// トークンの持ち主から指定したユーザーへいいね！を送信します。
	addLike := entities.UserLike{
		UserID:    id,
		PartnerID: partnerID,
	}
	err = likeRepo.Create(addLike)
	if err != nil {
		return postLikeInternalServerErrorResponse()
	}

	// いいね！を送信したお相手が自分にいいね！をしていたかどうかを検証します。
	liked, err = likeRepo.GetLikeBySenderIDReceiverID(addLike.PartnerID, addLike.UserID)
	if err != nil {
		return postLikeInternalServerErrorResponse()
	}
	// 過去にいいね！されていたらマッチングします。
	if liked != nil {
		addMatch := entities.UserMatch{
			UserID:    liked.UserID,
			PartnerID: liked.PartnerID,
		}
		matching := matchRepo.Create(addMatch)
		if matching != nil {
			return postLikeInternalServerErrorResponse()
		}

		return postLikeOK()
	}

	return postLikeOK()
}
