package userlike

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	/*
	1.	tokenのvalidation
	2.	tokenからuseridを取得
	3.	userIDからマッチ済みの相手matchIDを取得
	4.	useridからマッチ済み以外のいいねの受信リストを取得
	5.	いいねの受信リストからユーザーのプロフィールのリストを取得
	*/


	// Tokenがあるかどうか
	if p.Token == "" {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "No Token",
			})
	}

	// tokenからuserIDを取得

	rToken := repositories.NewUserTokenRepository()
	entToken, errToken := rToken.GetByToken(p.Token)
	if errToken != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error(tokenからuserIDを取得)",
			})
	}

	if entToken == nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized Token",
			})
	}

	sEntToken := entToken.Build()


	// matchIDsの取得

	rMatch := repositories.NewUserMatchRepository()
	matchIDs, errMatch := rMatch.FindAllByUserID(sEntToken.UserID)

	if errMatch != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}



	//fmt.Println("matchIDs",matchIDs)
	// マッチ済み以外のいいね受信リストを取得する
	rLike := repositories.NewUserLikeRepository()
	limit := int(p.Limit)
	if limit > 20 {
		limit = 20
	}
	offset := int(p.Offset)

	//fmt.Println("sEntToken.UserID",sEntToken.UserID)
	//fmt.Println("limit",limit)
	//fmt.Println("offset",offset)
	likes, errLike := rLike.FindGotLikeWithLimitOffset(sEntToken.UserID, limit, offset, matchIDs)
	if errLike != nil {
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	//fmt.Println("likes",likes)
	userLikes := entities.UserLikes(likes)

	//fmt.Println("userLikes",userLikes)
	sUsers := userLikes.Build() // userID partnerID createdAt UpdatedAtのリスト
	//fmt.Println("sUsers",sUsers)

	// sUsersが全て同じになってしまう。ポインタについては解決方法がわからない。

	rUser := repositories.NewUserRepository()

	// 上で取得した全てのpartnerIDについて、プロフィール情報を取得してpayloadsに格納する。
	var payloads []*models.LikeUserResponse
	for _, sUser := range sUsers{
		has, err := rUser.GetByUserID(sUser.UserID)
		
		if err != nil{
			return si.NewGetLikesInternalServerError().WithPayload(
				&si.GetLikesInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
		//entities.User -> models.LikeUserResponse
		r := models.LikeUserResponse{}
		r.ID = has.ID
		r.Nickname = has.Nickname
		r.Tweet = has.Tweet
		r.Introduction = has.Introduction
		r.ResidenceState = has.ResidenceState
		r.HomeState = has.HomeState
		r.Education = has.Education
		r.Job = has.Job
		r.AnnualIncome = has.AnnualIncome
		r.Height = has.Height
		r.BodyBuild = has.BodyBuild
		r.MaritalStatus = has.MaritalStatus
		r.Child = has.Child
		r.WhenMarry = has.WhenMarry
		r.WantChild = has.WantChild
		r.Smoking = has.Smoking
		r.Drinking = has.Drinking
		r.Holiday = has.Holiday
		r.HowToMeet = has.HowToMeet
		r.CostOfDate = has.CostOfDate
		r.NthChild = has.NthChild
		r.Housework = has.Housework
		r.ImageURI = has.ImageURI
		r.CreatedAt = has.CreatedAt
		r.UpdatedAt = has.UpdatedAt
		/* r.LikedAt = (探しても見つからない)*/
		payloads = append(payloads,&r)
	}

	//fmt.Println("payloads",payloads)
	return si.NewGetLikesOK().WithPayload(payloads)
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	/*
	1.	Tokenのバリデーション
	2.	tokenから送信者のuseridを取得
	3.	送信者のuseridから送信者のプロフィルを持ってきて性別を確認
	4.	p.useridから送信相手のプロフィルを持ってきて異性かどうか確認
	5.	すでにいいねしているか確認
	6.	いいねを送信
	*/


	// Tokenがあるかどうか
	if p.Params.Token == "" {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "401",
				Message: "No Token",
			})
	}

	// tokenから送信者のuserIDを取得

	rToken := repositories.NewUserTokenRepository()
	entToken, errToken := rToken.GetByToken(p.Params.Token)
	if errToken != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entToken == nil {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized Token",
			})
	}

	sEntToken := entToken.Build()


	// 送信者のuseridから送信者のプロフィルを持ってきて性別を確認
	// genderを確認するためだけに、useridからプロフィルを取得する……
	rUser := repositories.NewUserRepository()
	entUser, errUser := rUser.GetByUserID(sEntToken.UserID)
	if errUser != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entUser == nil { // entUserがnilになることはないはずだが、一応書いておく
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	gender := entUser.GetOppositeGender()

	// 送信相手のuseridから送信相手のプロフィルを持ってきて性別を確認
	// genderを確認するためだけに、useridからプロフィルを取得する……

	// userを設定する
	entUser2, errUser2 := rUser.GetByUserID(p.UserID)
	if errUser2 != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entUser2 == nil { // entUserがnilになることはないはずだが、一応書いておく
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 異性かどうかの確認
	if entUser2.Gender != gender {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}


	// すでにいいねしているかどうか確認する
	// userIDはいいねを送った人, partnerIDはいいねを受け取った人
	rLike := repositories.NewUserLikeRepository()
	entLike, errLike := rLike.GetLikeBySenderIDReceiverID(sEntToken.UserID, p.UserID)
	if errLike != nil {
		return si.NewPostLikeInternalServerError().WithPayload(
		&si.PostLikeInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
	}
	// すでにいいねしている場合
	if entLike != nil {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	var userLike entities.UserLike
	userLike.UserID = sEntToken.UserID
	userLike.PartnerID = p.UserID
	userLike.CreatedAt = strfmt.DateTime(time.Now())
	userLike.UpdatedAt = userLike.CreatedAt
	// いいねを送信する
	errLikeCreate := rLike.Create(userLike)
	if errLikeCreate != nil {
		//fmt.Println(errLikeCreate)
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
