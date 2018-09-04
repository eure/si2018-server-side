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

	// ユーザーID取得用
	ut, _ := t.GetByToken(p.Token)
	// match済みのユーザーを取得
	match, _ := m.FindAllByUserID(ut.UserID)
	// 自分が既にマッチングしている全てのお相手のUserIDを返す
	like, _ := l.FindGetLikeWithLimitOffset(ut.UserID, int(p.Limit), int(p.Offset), match)

	// 明示的に型宣言
	var lu entities.LikeUserResponses
	// ApplyUserの形にしてBuildしないといけないので
	// ApplyUserの型に変換する
	for _, likeUser := range like {
		// 構造体の初期化
		r := entities.LikeUserResponse{}
		user, _ := u.GetByUserID(likeUser.UserID)
		r.ApplyUser(*user)
		lu = append(lu, r)
	}
	sEnt := lu.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	l := repositories.NewUserLikeRepository()
	u := repositories.NewUserRepository()

	// ユーザーID取得用
	ut, _ := t.GetByToken(p.Params.Token)
	// ユーザーがすでにいいねしている人を取得
	all, _ := l.FindLikeOnley(ut.UserID)
	// 利用しているユーザーといいねしたいユーザー
	user, _ := u.GetByUserID(ut.UserID)
	likeUser, _ := u.GetByUserID(p.UserID)

	// 性別の確認
	if CheckGenderUserID(likeUser, user) {
		// すでにいいねしているかの確認
		if CheckLikeUserID(all, p.UserID) {
			// いいねの初期値の定義
			like := entities.UserLike{
				UserID:    ut.UserID,
				PartnerID: p.UserID,
				CreatedAt: strfmt.DateTime(time.Now()),
				UpdatedAt: strfmt.DateTime(time.Now()),
			}
			// いいねをインサート
			err := l.Create(like)
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
	// マッチングするかの確認
	err := LikeMatch(p.UserID, ut.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
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

// お互いいいねになったらマッチングさせる
func LikeMatch(likeUserID, userID int64) error {
	var err error
	m := repositories.NewUserMatchRepository()
	l := repositories.NewUserLikeRepository()

	r, _ := l.GetLikeBySenderIDReceiverID(likeUserID, userID)
	if r != nil {
		// マッチングの初期値を定義
		match := entities.UserMatch{
			UserID:    userID,
			PartnerID: likeUserID,
			CreatedAt: strfmt.DateTime(time.Now()),
			UpdatedAt: strfmt.DateTime(time.Now()),
		}
		// マッチングのインサート
		err = m.Create(match)
	}
	return err
}
