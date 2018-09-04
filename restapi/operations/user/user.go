package user

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/k0kubun/pp"
	"fmt"
	"github.com/eure/si2018-server-side/entities"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	/*
	1. tokenのvalidation
	2. tokenからuseridを取得
	3. useridからいいねを送受信した人を取得
	4. useridから異性がどちらであるかを判断
	5. いいねを送受信した人以外で異性の人のリストを取得する
	*/


	// Tokenがあるかどうか
	if p.Token == "" {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "No Token",
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
	ids, errLike := rLike .FindLikeAll(sEntToken.UserID)
	if errLike != nil{
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
	limit := int(p.Limit)
	if limit > 20 {
		limit = 20 // 上限を20までにしました
	}
	offset := int(p.Offset)
	// 異性のリストの取得
	usersFind, errFind := rUser.FindWithCondition(limit, offset, gender, ids)
	if errFind != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}


	var users entities.Users


	for _, user := range usersFind{
		users = append(users, user)
	}

	sUsers := users.Build()

	return si.NewGetUsersOK().WithPayload(sUsers)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {

	// Tokenのチェック
	if p.Token == "" {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "No Token",
			})
	}

	r1 := repositories.NewUserTokenRepository()
	ent1, err1 := r1.GetByToken(p.Token)
	if err1 != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if ent1 == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "401",
				Message: "Unauthorized Token",
			})
	}

	// Tokenがあった場合の処理
	r := repositories.NewUserRepository()

	ent, err := r.GetByUserID(p.UserID)
	fmt.Println(err)


	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	if ent == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found",
			})
	}

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	// 1.tokenが正しいか
	// 2.tokenとUSERIDが対応しているか（もし対応していなかったらforbidden）
	// 3.更新

	// Tokenのチェック
	if p.Params.Token == "" {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "No Token",
			})
	}

	r1 := repositories.NewUserTokenRepository()
	ent1, err1 := r1.GetByToken(p.Params.Token)

	if err1 != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	// tokenに対応するuseridが見つからない時
	if ent1 == nil {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized Token",
			})
	}

	// tokenとUSERIDのチェック
	sEnt1 := ent1.Build()
	pp.Print(sEnt1)

	// 他人のプロフィルを更新しようとしていた場合
	if sEnt1.UserID != p.UserID {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code:	"403",
				Message:"Forbidden",
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
	u.ImageURI = pa.ImageURI

	r := repositories.NewUserRepository()
	err2 := r.Update(&u)
	if err2 != nil  {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	ent, err := r.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	sEnt := ent.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
	/*疑問点……err, entの変数名をいちいち変えるべきか？*/
}
