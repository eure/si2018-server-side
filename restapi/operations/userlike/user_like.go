package userlike

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	// トークンからユーザーIDを取得
	ruserR := repositories.NewUserRepository()
	user, _ := ruserR.GetByToken(p.Token)
	userid := user.ID

	// マッチング済みのユーザーを取得
	ruserM := repositories.NewUserMatchRepository()
	matchedids, _ := ruserM.FindAllByUserID(userid)

	r := repositories.NewUserLikeRepository()

	// いいねされたやつらを集める
	var ent entities.UserLikes
	ent, _ = r.FindGotLikeWithLimitOffset(userid,int(p.Limit),int(p.Offset),matchedids)

	// applied メソッドによって変換されたUser's'がほしい。
	// とりあえずほしいから，格納先を用意してあげる。
	var appliedusers entities.LikeUserResponses
	for _,m := range ent {
		var applied = entities.LikeUserResponse{}
		likeduser , _ := ruserR.GetByUserID(m.UserID)
		applied.ApplyUser(*likeduser)
		appliedusers = append (appliedusers, applied)
	}

	// aplyされた結果がbuildされればいい
	sEnt := appliedusers.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	return si.NewPostLikeOK()
}
