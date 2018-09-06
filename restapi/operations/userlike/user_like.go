package userlike

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"time"
)


//- GET {hostname}/api/1.0/likes
//- 自分にいいね！を送ってくれたユーザーを20件レスポンスとして返してください
//- ページネーションを実装してください
//- TokenのValidation処理を実装してください
func GetLikes(p si.GetLikesParams) middleware.Responder {
	repUserLike := repositories.NewUserLikeRepository()
	repUserToken := repositories.NewUserTokenRepository() // tokenからユーザーのIDを取得するため
	repUserMatch := repositories.NewUserMatchRepository() // ユーザーとマッチングしているユーザーを取得するため
	repUser := repositories.NewUserRepository() // idからユーザーを取得するため

	// Tokenのバリデーション
	err := repUserToken.ValidateToken(p.Token)
	if err != nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code: "401",
				Message: "Token Is Invalid",
			})
	}

	// bad Request
	if (p.Limit < 1) || (p.Offset < 0) {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// トークンからユーザーのIDを取得
	loginUserToken, _ := repUserToken.GetByToken(p.Token)

	// そのユーザーとマッチングしているユーザーのIDを全取得
	matchUserIDs, err := repUserMatch.FindAllByUserID(loginUserToken.UserID) 
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// マッチングしている人を除く、ユーザーをいいねした人を全取得
	userLikesForMe, err := repUserLike.FindGotLikeWithLimitOffset(loginUserToken.UserID, int(p.Limit), int(p.Offset), matchUserIDs)
	if err != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	//if entsUserLike_LikedMe == nil {
	//	return si.NewGetLikesOK().WithPayload(nil)
	//}

	// ユーザーをいいねした人のIDの配列を取得
	var likedUserIDs []int64
	for _, userLikeForMe := range userLikesForMe {
		likedUserIDs = append(likedUserIDs, userLikeForMe.UserID)
	}

	// idの配列からユーザーを取得
	likedMeUsers , _:= repUser.FindByIDs(likedUserIDs)

	// return用のmodelを作るためのLikeUserResponsesのエンティティ
	var likeUserReses entities.LikeUserResponses
	//
	for i, likedMeUsers := range likedMeUsers {
		var likeUserRese = entities.LikeUserResponse{}
		likeUserRese.ApplyUser(likedMeUsers)
		likeUserRese.LikedAt = userLikesForMe[i].CreatedAt
		likeUserReses = append(likeUserReses, likeUserRese)
	}

	likeUserRese := likeUserReses.Build()

	return si.NewGetLikesOK().WithPayload(likeUserRese)
}


//- POST {hostname}/api/1.0/likes/{userID}
//- 相手にいいね！を送信してください
//- TokenのValidation処理を実装してください
//- ※ お互いいいね！を送るとmatchingになります
//- ※ 同性へはいいね！ができません
//- ※ 同じ人へは2回いいね！ができません
func PostLike(p si.PostLikeParams) middleware.Responder {
	repUserLike := repositories.NewUserLikeRepository()
	repUserToken := repositories.NewUserTokenRepository() //tokenからログインユーザーのIDを取るため
	repUser := repositories.NewUserRepository() // ログインユーザー、いいねしたユーザーを取得するため
	repUserMatch := repositories.NewUserMatchRepository() // お互いいいねした時にマッチングのため

	// tokenのバリデーション
	err := repUserToken.ValidateToken(p.Params.Token)
	if err != nil {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code: "401",
				Message: "Token Is Invalid",
			})
	}

	// Bad Request
	if p.UserID < 0 {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// tokenからユーザーIDを取得し、そのIDのユーザーを取得
	loginUserToken, _ := repUserToken.GetByToken(p.Params.Token)
	loginUser, err := repUser.GetByUserID(loginUserToken.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// パラメータのIDからいいねするユーザーを取得
	likeUser, err := repUser.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 同性時のエラーハンドリング
	if loginUser.Gender == likeUser.Gender {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// 同じ人への２回いいねのエラーハンドリング
	// FindILikedAllで自分がいいねした人のIDを返す
	alreadyLikedIDs, err := repUserLike.FindILikedAll(loginUser.ID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	for _, alreadyLikedID := range alreadyLikedIDs {
		if likeUser.ID == alreadyLikedID {
			return si.NewPostLikeBadRequest().WithPayload(
				&si.PostLikeBadRequestBody{
					Code:    "400",
					Message: "Bad Request",
				})
		}
	}

	// UserLikeの作成
	var userLike entities.UserLike
	userLike.UserID = loginUser.ID
	userLike.PartnerID = likeUser.ID
	userLike.CreatedAt = strfmt.DateTime(time.Now())
	userLike.UpdatedAt = userLike.CreatedAt

	err = repUserLike.Create(userLike)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// いいねした人がいいねした人のIDを全取得
	likedIDs, err := repUserLike.FindILikedAll(likeUser.ID)
	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// いいねした人が自分をいいねしているか確認し、いいねしていた場合マッチングさせる
	for _, likedID := range likedIDs {
		if loginUser.ID == likedID {
			var userMatch entities.UserMatch
			userMatch.UserID = userLike.PartnerID
			userMatch.PartnerID = userLike.UserID
			userMatch.CreatedAt = userLike.CreatedAt
			userMatch.UpdatedAt = userLike.UpdatedAt

			err := repUserMatch.Create(userMatch)
			if err != nil {
				return si.NewPostLikeInternalServerError().WithPayload(
					&si.PostLikeInternalServerErrorBody{
						Code:    "500",
						Message: "Internal Server Error",
					})
			}
		}
	}

	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code: "200",
			Message: "OK",
		})
}
