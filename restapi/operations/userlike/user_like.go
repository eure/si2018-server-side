package userlike

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/runtime/middleware"
	
	"time"
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
	r := repositories.NewUserLikeRepository()
	u := repositories.NewUserRepository()
	t := repositories.NewUserTokenRepository()
	m := repositories.NewUserMatchRepository()

	// p.Params からいいねの送信ユーザーのUserID,Userを取得
	FromUserParam := p.Params
	FromUserToken, _ := t.GetByToken(FromUserParam.Token)
	FromUserID := FromUserToken.UserID
	FromUser, _ := u.GetByToken(FromUserToken.Token)

	// p.UserID からいいねの受信ユーザーのUserID,Userを取得
	ToUserID := p.UserID
	ToUser, _ := u.GetByUserID(ToUserID)
	
	// 同性かどうか
	if ToUser.Gender != FromUser.Gender {
		// いいねをすでに送信しているか
		SendLikeResult, err1 := r.GetLikeBySenderIDReceiverID(FromUserID, ToUserID)
		if err1 != nil {
			return outPutPostStatus(500)
		} else if SendLikeResult != nil {
			return outPutPostStatus(400)
		} else {
			// いいねを作成
			var Like entities.UserLike
			Like.UserID = FromUserID
			Like.PartnerID = ToUserID
			Like.CreatedAt = strfmt.DateTime(time.Now())
			Like.UpdatedAt = strfmt.DateTime(time.Now())

			// いいねを送信
			err2 := r.Create(Like)
			if err2 != nil {
				return outPutPostStatus(500)
			}
			// 相手もいいねしてるかどうか (1件だけ)
			SendLikeResult, err1 = r.GetLikeBySenderIDReceiverID(ToUserID, FromUserID)
			if SendLikeResult != nil {
				// マッチングを作成
				var NewMatching entities.UserMatch
				NewMatching.UserID = FromUserID
				NewMatching.PartnerID = ToUserID
				NewMatching.CreatedAt = strfmt.DateTime(time.Now())
				NewMatching.UpdatedAt = strfmt.DateTime(time.Now())
				_ = m.Create(NewMatching)
			}
		}
	} else {
		/*return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
		*/
		return outPutPostStatus(500)
	}

	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:   "200",
			Message: "OK",
		})

}

func outPutPostStatus (num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	case 400:
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	return nil
}
