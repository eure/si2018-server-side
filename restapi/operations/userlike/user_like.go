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
	// 自分にいいねをした人のidをとってくる
	// 自分のtokenからユーザーのエンティティのインスタンスを作成
	// マッチングしている人のidを取得
	// FindGotLikeWithLimitOffsetでいいねをした人をとってくる
	// エンティティのlikeUserをlikeUserResponseにキャストする

	////repUserLike := repositories.NewUserLikeRepository()
	//repUserToken :=  repositories.NewUserTokenRepository() //tokenからログインユーザーを取得するため
	//repUserMatch := repositories.NewUserMatchRepository() // すでにマッチした人のidを取得するため
	////repUser := repositories.NewUserRepository() // 自分をいいねしたユーザーを取得するため
	//
	//// トークンからログインユーザー取得
	//entUserToken, _ := repUserToken.GetByToken(p.Token)
	//
	//// 自分とマッチングしているユーザーのIDを取得
	//matchUserIDs, _ := repUserMatch.FindAllByUserID(entUserToken.UserID)
	//
	//for i, id := range matchUserIDs {
	//	fmt.Print(id)
	//	fmt.Print(i)
	//}

	//
	////マッチングしている人を除いて、自分にいいねをした人を取得
	//entsUserLike, _ := repUserLike.FindGotLikeWithLimitOffset(entUser.UserID, int(p.Limit), int(p.Offset), matchIDs)
	//entUserLikes := entities.UserLikes(entsUserLike) // Userlikeの配列をUserlikesにキャスト

	//// 自分をLikeしたユーザーのIDを格納するための配列
	//userLikeMeIDs := make([]int64, 0, len(entUserLikes))
	//
	//// 自分をLikeしたユーザーのIDを格納
	//for _, entUserLike := range entUserLikes {
	//	userLikeMeIDs = append(userLikeMeIDs, entUserLike.PartnerID)
	//
	//}
	//
	//println(userLikeMeIDs)
	//
	//// idをもとにユーザーの配列を取得
	//entUsersLikeMe, _ := repUser.FindByIDs(userLikeMeIDs)
	//
	////Responseの配列を作成
	//entUserLikeReses := make(entities.LikeUserResponses, 0, len(entUserLikes))
	//
	////ユーザの配列からresponseを作成
	//for i, entUserLikeMe := range entUsersLikeMe {
	//	entUserLikeReses[i].ApplyUser(entUserLikeMe)
	//	entUserLikeReses[i].LikedAt = entUserLikes[i].UpdatedAt
	//}
	//
	//userLikeReses := entUserLikeReses.Build()

	//return si.NewGetLikesOK().WithPayload(userLikeReses)
	return si.NewGetLikesOK()
}


//- POST {hostname}/api/1.0/likes/{userID}
////- 相手にいいね！を送信してください
////- TokenのValidation処理を実装してください
////- ※ お互いいいね！を送るとmatchingになります
////- ※ 同性へはいいね！ができません
////- ※ 同じ人へは2回いいね！ができません
func PostLike(p si.PostLikeParams) middleware.Responder {
	var entUserLike entities.UserLike

	repUserLike := repositories.NewUserLikeRepository()
	repUserToken := repositories.NewUserTokenRepository() //tokenからログインユーザーのIDを取るため
	repUser := repositories.NewUserRepository() // ログインユーザー、いいねしたユーザーを取得するため
	repUserMatch := repositories.NewUserMatchRepository() // お互いいいねした時にマッチングのため

	// tokenからユーザーIDを取得し、そのIDのユーザーを取得
	entLoginUserToken, _ := repUserToken.GetByToken(p.Params.Token)
	entLoginUser, _ := repUser.GetByUserID(entLoginUserToken.UserID)

	// パラメータのIDからいいねするユーザーを取得
	entLikedUser, _ := repUser.GetByUserID(p.UserID)

	// 同性時のエラーハンドリング
	if entLoginUser.Gender == entLikedUser.Gender {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "404",
				Message: "I`m sorry, This service not support gay.",
			})
	}

	// 同じ人への２回いいねのエラーハンドリング
	// FindILikedAllで自分がいいねした人のIDを返す
	exceptIDs, _ := repUserLike.FindILikedAll(entLoginUser.ID)
	for _, exceptID := range exceptIDs {
		if entLikedUser.ID == exceptID {
			return si.NewPostLikeBadRequest().WithPayload(
				&si.PostLikeBadRequestBody{
					Code:    "404",
					Message: "You already liked this user.",
				})
		}
	}

	entUserLike.UserID = entLoginUser.ID
	entUserLike.PartnerID = entLikedUser.ID
	entUserLike.CreatedAt = strfmt.DateTime(time.Now())
	entUserLike.UpdatedAt = entUserLike.CreatedAt

	err := repUserLike.Create(entUserLike)

	if err != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// いいねした相手が自分をいいねしていた場合、マッチングさせる
	//　いいねした人のいいねに自分がいるか確認
	LikedIDs, _ := repUserLike.FindILikedAll(entLikedUser.ID)
	for _, LikedID := range LikedIDs {
		if entLoginUser.ID == LikedID {
			var entUserMatch entities.UserMatch
			// マッチさせる
			entUserMatch.UserID = entUserLike.PartnerID
			entUserMatch.PartnerID = entUserLike.UserID
			entUserMatch.CreatedAt = entUserLike.CreatedAt
			entUserMatch.UpdatedAt = entUserLike.UpdatedAt

			err := repUserMatch.Create(entUserMatch)
			if err != nil {
				return si.NewPostLikeInternalServerError().WithPayload(
					&si.PostLikeInternalServerErrorBody{
						Code:    "500",
						Message: "Internal Server Error",
					})
			}

			return si.NewPostLikeOK().WithPayload(
				&si.PostLikeOKBody{
					Code: "201",
					Message: "You matched",
				})

		}
	}

	return si.NewPostLikeOK().WithPayload(
		&si.PostLikeOKBody{
			Code: "201",
			Message: "Good Luck",
		})
}
