package userlike

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	
	"time"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	m := repositories.NewUserMatchRepository()
	u := repositories.NewUserRepository()
	r := repositories.NewUserLikeRepository()
	
	// loginUserのUser entitiesを取得 (Validation)
	loginUserToken := p.Token
	loginUser,err := t.GetByToken(loginUserToken)
	if err != nil {
		return outPutGetStatus(500)
	}
	if loginUser == nil {
		return outPutGetStatus(401)
	}

	// limit が20かどうか検出
	if p.Limit != int64(20) {
		return outPutGetStatus(400)
	}
	// offset が0以上かどうか検出
	if p.Offset < int64(0) {
		return outPutGetStatus(400)
	}
	
	loginUserID := loginUser.UserID

	// マッチング済みのユーザーを取得
	matchedIDs, err := m.FindAllByUserID(loginUserID)
	if err != nil {
		return outPutGetStatus(500)
	}
	
	// いいねされたやつらを集める
	var ent entities.UserLikes
	ent, err = r.FindGotLikeWithLimitOffset(loginUserID,int(p.Limit),int(p.Offset),matchedIDs)
	if err != nil {
		return outPutGetStatus(500)
	}

	// applied メソッドによって変換されたUser's'がほしい。
	// とりあえずほしいから，格納先を用意してあげる。
	var appliedUsers entities.LikeUserResponses
	for _,m := range ent {
		var applied = entities.LikeUserResponse{}
		likedUser , _ := u.GetByUserID(m.UserID)
		applied.ApplyUser(*likedUser)
		appliedUsers = append (appliedUsers, applied)
	}

	// aplyされた結果がbuildされればいい
	sEnt := appliedUsers.Build()
	return si.NewGetLikesOK().WithPayload(sEnt)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	r := repositories.NewUserLikeRepository()
	u := repositories.NewUserRepository()
	t := repositories.NewUserTokenRepository()
	m := repositories.NewUserMatchRepository()

	// fromUser == loginUser
	// tokenから User entitiesを取得 (Validation)
	token := p.Params.Token
	fromUserToken, _ := t.GetByToken(token)
	fromUserID := fromUserToken.UserID
	fromUser, _ := u.GetByUserID(fromUserID)
	
	// p.UserID からいいねの受信ユーザーのUserID,Userを取得
	toUserID := p.UserID
	toUser, _ := u.GetByUserID(toUserID)
	
	// 自分に送ろうとしていないか
	if fromUserID == toUserID {
		return outPutPostStatus(400)
	}
	
	// 同性かどうか
	if toUser.Gender != fromUser.Gender {
		// いいねをすでに送信しているか
		SendLikeResult, err := r.GetLikeBySenderIDReceiverID(fromUserID, toUserID)
		if err != nil {
			return outPutPostStatus(500)
		} else if SendLikeResult != nil {
			return outPutPostStatus(400)
		} else {
			// いいねを作成
			var Like entities.UserLike
			Like.UserID = fromUserID
			Like.PartnerID = toUserID
			Like.CreatedAt = strfmt.DateTime(time.Now())
			Like.UpdatedAt = strfmt.DateTime(time.Now())

			// いいねを送信
			err := r.Create(Like)
			if err != nil {
				return outPutPostStatus(500)
			}
			// 相手もいいねしてるかどうか (1件だけ)
			SendLikeResult, err = r.GetLikeBySenderIDReceiverID(toUserID, fromUserID)
			if SendLikeResult != nil {
				// マッチングを作成
				var NewMatching entities.UserMatch
				NewMatching.UserID = fromUserID
				NewMatching.PartnerID = toUserID
				NewMatching.CreatedAt = strfmt.DateTime(time.Now())
				NewMatching.UpdatedAt = strfmt.DateTime(time.Now())
				_ = m.Create(NewMatching)
				
			}
		}
	} else {
		return outPutPostStatus(400)
	}

	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code:   "200",
			Message: "OK",
		})

}

func outPutGetStatus (num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	case 401:
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	case 400:
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	return nil
}

func outPutPostStatus (num int) middleware.Responder {
	switch num {
	case 500:
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	case 401:
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized (トークン認証に失敗)",
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
