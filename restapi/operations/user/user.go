package user

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	/*
		1. tokenのvalidation
		2. tokenからuseridを取得
		3. useridからいいねを送受信した人を取得
		4. useridから異性がどちらであるかを判断
		5. いいねを送受信した人以外で異性の人のリストを取得する
		6. 画像の取得
	*/

	// Tokenがあるかどうか
	if p.Token == "" {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Required",
			})
	}

	// tokenからuserIDを取得

	rToken := repositories.NewUserTokenRepository()
	entToken, errToken := rToken.GetByToken(p.Token)
	if errToken != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entToken == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized Token",
			})
	}

	sEntToken := entToken.Build()

	// useridからすでにいいねを送受信した人のリストを取得
	rLike := repositories.NewUserLikeRepository()
	ids, errLike := rLike.FindLikeAll(sEntToken.UserID)
	if errLike != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// genderを設定するためだけに、useridからプロフィルを取得する……
	rUser := repositories.NewUserRepository()
	entUser, errUser := rUser.GetByUserID(sEntToken.UserID)
	if errUser != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entUser == nil { // entUserがnilになることはないはずだが、一応書いておく
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	gender := entUser.GetOppositeGender()

	// p.Limit, Offsetがなぜかint64なので、intに変換しないといけない。
	// limitが0だと全件取得。
	limit := int(p.Limit)
	offset := int(p.Offset)
	if limit <= 0 || offset < 0 {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}
	// 異性のリストの取得
	usersFind, errFind := rUser.FindWithCondition(limit, offset, gender, ids)
	if errFind != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	users := entities.Users(usersFind)

	var userIdList []int64
	// usersからidのリストを取得
	for _, user := range users {
		userIdList = append(userIdList, user.ID)
	}

	if userIdList == nil {
		var sUsers []*models.User
		return si.NewGetUsersOK().WithPayload(sUsers)
	}

	// 画像取得
	rImage := repositories.NewUserImageRepository()
	entImages, errImages := rImage.GetByUserIDs(userIdList)
	if errImages != nil || entImages == nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// id -- pathの対応リストを作成
	idPaths := map[int64]string{}
	for _, entImage := range entImages {
		idPaths[entImage.UserID] = entImage.Path
	}
	sUsers := users.Build()

	// 画像のpathを結合
	for _, sUser := range sUsers {
		sUser.ImageURI = idPaths[sUser.ID]
	}
	return si.NewGetUsersOK().WithPayload(sUsers)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	/*
		1.	tokenのvalidation
		2.	tokenからuseridを取得
		3.	profileを取得
		4.	画像URIを取得
	*/

	// Tokenのチェック
	if p.Token == "" {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Required",
			})
	}

	rToken := repositories.NewUserTokenRepository()
	entToken, errToken := rToken.GetByToken(p.Token)
	if errToken != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entToken == nil {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// Tokenがあった場合の処理
	rUser := repositories.NewUserRepository()

	entUser, errUser := rUser.GetByUserID(p.UserID)

	if errUser != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if entUser == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found",
			})
	}

	// 画像取得

	rImage := repositories.NewUserImageRepository()
	entImage, errImage := rImage.GetByUserID(p.UserID)
	if errImage != nil || entImage == nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	entUser.ImageURI = entImage.Path

	sEnt := entUser.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	/*
		1.tokenが正しいか
		2.tokenとUSERIDが対応しているか（もし対応していなかったらforbidden）
		3.更新
	*/

	// Tokenのチェック
	if p.Params.Token == "" {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Required",
			})
	}

	rToken := repositories.NewUserTokenRepository()
	entToken, errToken := rToken.GetByToken(p.Params.Token)

	if errToken != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	// tokenに対応するuseridが見つからない時
	if entToken == nil {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// tokenとUSERIDのチェック
	sEnt1 := entToken.Build()

	// 他人のプロフィルを更新しようとしていた場合
	if sEnt1.UserID != p.UserID {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code:    "403",
				Message: "Forbidden",
			})
	}

	pa := p.Params
	u := entities.User{}

	u.ID = p.UserID
	u.Nickname = pa.Nickname
	u.Tweet = pa.Tweet
	u.Introduction = pa.Introduction
	u.ResidenceState = pa.ResidenceState
	u.HomeState = pa.HomeState
	u.Education = pa.Education
	u.Job = pa.Job
	u.AnnualIncome = pa.AnnualIncome
	u.Height = pa.Height
	u.BodyBuild = pa.BodyBuild
	u.MaritalStatus = pa.MaritalStatus
	u.Child = pa.Child
	u.WhenMarry = pa.WhenMarry
	u.WantChild = pa.WantChild
	u.Smoking = pa.Smoking
	u.Drinking = pa.Drinking
	u.Holiday = pa.Holiday
	u.HowToMeet = pa.HowToMeet
	u.CostOfDate = pa.CostOfDate
	u.NthChild = pa.NthChild
	u.Housework = pa.Housework
	// u.ImageURI = pa.ImageURI

	rUser := repositories.NewUserRepository()
	errUpdate := rUser.Update(&u)
	if errUpdate != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	entUser, errUser := rUser.GetByUserID(p.UserID)
	if errUser != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	sEnt := entUser.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}
