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

// DB アクセス: 5 回
// 計算量: O(N)
func GetLikes(p si.GetLikesParams) middleware.Responder {
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Token)
		if err != nil {
			return si.GetLikesThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return si.GetLikesThrowUnauthorized("GetByToken failed")
		}
		id = token.UserID
	}
	// もらったいいねを取得
	var like []entities.UserLike
	{
		likeRepo := repositories.NewUserLikeRepository()
		matchRepo := repositories.NewUserMatchRepository()
		matched, err := matchRepo.FindAllByUserID(id)
		if err != nil {
			return si.GetLikesThrowInternalServerError("FindAllByUserID", err)
		}
		like, err = likeRepo.FindGotLikeWithLimitOffset(id, int(p.Limit), int(p.Offset), matched)
		if err != nil {
			return si.GetLikesThrowInternalServerError("FindGotLikeWithLimitOffset", err)
		}
	}
	// 相手の ID の取得
	ids := make([]int64, 0)
	for _, l := range like {
		ids = append(ids, l.UserID)
	}
	count := len(like)
	mapping := make(map[int64]int)
	for i, l := range like {
		mapping[l.UserID] = i
	}
	// いいねに紐づくユーザー情報を取得
	users := make([]entities.User, count)
	{
		userRepo := repositories.NewUserRepository()
		shuffledUsers, err := userRepo.FindByIDs(ids)
		if err != nil {
			return si.GetLikesThrowInternalServerError("FindByIDs", err)
		}
		if len(shuffledUsers) != count {
			return si.GetLikesThrowBadRequest("FindByIDs failed")
		}
		for _, u := range shuffledUsers {
			users[mapping[u.ID]] = u
		}
	}
	// 対応する画像の取得
	images := make([]entities.UserImage, count)
	{
		imageRepo := repositories.NewUserImageRepository()
		shuffledImages, err := imageRepo.GetByUserIDs(ids)
		if err != nil {
			return si.GetLikesThrowInternalServerError("GetByUserIDs", err)
		}
		if len(shuffledImages) != count {
			return si.GetLikesThrowBadRequest("GetByUserIDs failed")
		}
		for _, im := range shuffledImages {
			images[mapping[im.UserID]] = im
		}
	}
	// 以上の情報をまとめる
	sEnt := make([]*models.LikeUserResponse, 0)
	for i, l := range like {
		response := entities.LikeUserResponse{LikedAt: l.UpdatedAt}
		response.ApplyUser(users[i])
		swaggerLike := response.Build()
		swaggerLike.ImageURI = images[i].Path
		sEnt = append(sEnt, &swaggerLike)
	}

	return si.NewGetLikesOK().WithPayload(sEnt)
}

// DB アクセス: 6,7 回
// 計算量: O(1)
func PostLike(p si.PostLikeParams) middleware.Responder {
	if p.Params.Token == "" {
		return si.PostLikeThrowBadRequest("missing token")
	}
	userRepo := repositories.NewUserRepository()
	likeRepo := repositories.NewUserLikeRepository()
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Params.Token)
		if err != nil {
			return si.PostLikeThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return si.PostLikeThrowUnauthorized("GetByToken failed")
		}
		id = token.UserID
	}
	// 同性にいいねは送れないので, 性別情報を取得する
	var oppositeGender string
	{
		user, err := userRepo.GetByUserID(id)
		if err != nil {
			return si.PostLikeThrowInternalServerError("GetByUserID", err)
		}
		if user == nil {
			return si.PostLikeThrowBadRequest("GetByUserID failed")
		}
		oppositeGender = user.GetOppositeGender()
	}
	// 送る相手の情報を取得
	partner, err := userRepo.GetByUserID(p.UserID)
	if err != nil {
		return si.PostLikeThrowInternalServerError("GetByUserID", err)
	}
	if partner == nil {
		return si.PostLikeThrowBadRequest("GetByUserID failed")
	}
	if partner.Gender != oppositeGender {
		return si.PostLikeThrowBadRequest("genders must be opposite")
	}
	// いいねが重複していないかの確認
	like, err := likeRepo.GetLikeBySenderIDReceiverID(id, partner.ID)
	if err != nil {
		return si.PostLikeThrowInternalServerError("GetLikeBySenderIDReceiverID", err)
	}
	if like != nil {
		return si.PostLikeThrowBadRequest("like action duplicates")
	}
	// 相手からいいねが来ていたかの確認
	reverse, err := likeRepo.GetLikeBySenderIDReceiverID(partner.ID, id)
	if err != nil {
		return si.PostLikeThrowInternalServerError("GetLikeBySenderIDReceiverID", err)
	}
	now := strfmt.DateTime(time.Now())
	like = new(entities.UserLike)
	*like = entities.UserLike{
		UserID:    id,
		PartnerID: partner.ID,
		CreatedAt: now,
		UpdatedAt: now}
	// like を書き込んだあと, match を書き込むときにエラーが発生すると致命的
	// https://qiita.com/komattio/items/838ea5df68eb076e8099
	// transaction を利用してまとめて書きこむ必要がある
	err = likeRepo.Create(*like)
	if err != nil {
		return si.PostLikeThrowInternalServerError("Create", err)
	}
	// お互いにいいねしていればマッチング成立
	if reverse != nil {
		matchRepo := repositories.NewUserMatchRepository()
		match := entities.UserMatch{
			UserID:    partner.ID,
			PartnerID: id,
			CreatedAt: now,
			UpdatedAt: now}
		err = matchRepo.Create(match)
		if err != nil {
			return si.PostLikeThrowInternalServerError("Create", err)
		}
	}
	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:    "200",
			Message: "OK",
		})
}
