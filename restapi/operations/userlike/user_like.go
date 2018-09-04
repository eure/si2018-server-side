package userlike

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
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
	for _, likeUser := range like{
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
	return si.NewPostLikeOK()
}
