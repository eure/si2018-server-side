package userlike

import (
	"log"

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
	}
	if tokenOwner == nil {
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

	// ユーザーをいいね!しているお相手のIDの配列を取得します。
	ids := make([]int64, len(matches))
	for _, like := range likes {
		ids = append(ids, like.UserID)
	}

	usersChan := make(chan map[int64]entities.User, 1)
	imagesChan := make(chan map[int64]entities.UserImage, 1)
	errChan1 := make(chan error, 1)
	errChan2 := make(chan error, 1)

	go func(result chan error) {
		usrs, err := userRepo.FindByIDs(ids)
		if err != nil {
			log.Fatal(err)
			result <- err
		}

		mpU := map[int64]entities.User{}
		for _, u := range usrs {
			if u.ID > 0 {
				mpU[u.ID] = u
			}
		}

		usersChan <- mpU
		result <- nil

	}(errChan1)

	go func(result chan error) {
		imgs, err := imageRepo.GetByUserIDs(ids)
		if err != nil {
			log.Fatal(err)
			result <- err
		}

		mpI := map[int64]entities.UserImage{}
		for _, img := range imgs {
			if img.UserID > 0 {
				mpI[img.UserID] = img
			}
		}

		imagesChan <- mpI
		result <- nil

	}(errChan2)

	err = <-errChan1
	if err != nil {
		return getLikesInternalServerErrorResponse()
	}

	err = <-errChan2
	if err != nil {
		return getLikesInternalServerErrorResponse()
	}

	users := <-usersChan
	images := <-imagesChan

	var ents entities.LikeUserResponses
	for _, like := range likes {
		ent := entities.LikeUserResponse{}
		ent.ApplyUser(users[like.UserID])
		ent.ImageURI = images[like.UserID].Path
		ent.LikedAt = like.CreatedAt
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
