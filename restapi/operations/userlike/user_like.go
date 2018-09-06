package userlike

import (
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	l := repositories.NewUserLikeRepository()
	m := repositories.NewUserMatchRepository()
	u := repositories.NewUserRepository()

	// paramsの変数を定義
	paramsToken := p.Token
	paramsLimit := p.Limit
	paramsOffset := p.Offset

	// ユーザーID取得用
	token, err := t.GetByToken(paramsToken)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if token == nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// limitが20になっているかをvalidation
	if paramsLimit != int64(20) {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// match済みのユーザーを取得
	match, err := m.FindAllByUserID(token.UserID)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 自分が既にマッチングしている全てのお相手のUserIDを返す(limit,offset)
	like, err := l.FindGetLikeWithLimitOffset(token.UserID, int(paramsLimit), int(paramsOffset), match)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 明示的に型宣言
	var lu entities.LikeUserResponses
	var userIDs []int64

	for _, likeUser := range like {
		userIDs = append(userIDs, likeUser.UserID)
	}

	// いいねされているユーザーの情報をすべて取得する
	users, err := u.FindByIDs(userIDs)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// ApplyUserの形にしてBuildしないといけないので
	// ApplyUserの型に変換する
	for _, user := range users {
		// 構造体の初期化
		r := entities.LikeUserResponse{}
		// Userの情報を取得
		r.ApplyUser(user)
		lu = append(lu, r)
	}
	sEnt := lu.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	l := repositories.NewUserLikeRepository()
	u := repositories.NewUserRepository()
	m := repositories.NewUserMatchRepository()

	// paramsの変数を定義
	paramsToken := p.Params.Token
	paramsUserID := p.UserID

	// ユーザーID取得用
	token, err := t.GetByToken(paramsToken)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if token == nil {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// ユーザーがすでにいいねしている人を取得
	userIDs, err := l.FindLikeOnley(token.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 利用しているユーザーといいねしたいユーザー
	user, err := u.GetByUserID(token.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if user == nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	likeUser, err := u.GetByUserID(paramsUserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if likeUser == nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// 性別の確認
	if CheckGenderUserID(likeUser, user) {
		// すでにいいねしているかの確認
		if CheckLikeUserID(userIDs, paramsUserID) {
			// いいねの値の定義
			like := entities.UserLike{
				UserID:    token.UserID,
				PartnerID: paramsUserID,
				CreatedAt: strfmt.DateTime(time.Now()),
				UpdatedAt: strfmt.DateTime(time.Now()),
			}
			// いいねをインサート
			err = l.Create(like)
			if err != nil {
				return si.NewPostLikeInternalServerError().WithPayload(
					&si.PostLikeInternalServerErrorBody{
						Code:    "500",
						Message: "Internal Server Error",
					})
			}
		} else {
			return si.NewPostLikeBadRequest().WithPayload(
				&si.PostLikeBadRequestBody{
					Code:    "400",
					Message: "Bad Request",
				})
		}
	} else {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// お互いいいねになったらマッチングさせる
	// パートナー側がいいねをしているか確認してなければそのまま終了あればマッチングさせる
	r, err := l.GetLikeBySenderIDReceiverID(paramsUserID, token.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if r != nil {
		// マッチングの初期値を定義
		match := entities.UserMatch{
			UserID:    token.UserID,
			PartnerID: paramsUserID,
			CreatedAt: strfmt.DateTime(time.Now()),
			UpdatedAt: strfmt.DateTime(time.Now()),
		}
		// マッチングのインサート
		err = m.Create(match)
		if err != nil {
			return si.NewPostLikeInternalServerError().WithPayload(
				&si.PostLikeInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	}

	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:    "200",
			Message: "OK",
		})
}

// すでにいいねをしているか確認に使う関数
func CheckLikeUserID(likeID []int64, userID int64) bool {
	for _, l := range likeID {
		if l == userID {
			return false
		}
	}
	return true
}

// 性別が同じか確かめる
func CheckGenderUserID(likeUser, user *entities.User) bool {
	if likeUser.Gender == user.Gender {
		return false
	}
	return true
}
