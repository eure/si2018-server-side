package user

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	l := repositories.NewUserLikeRepository()
	u := repositories.NewUserRepository()

	// paramsの変数を定義
	paramsToken := p.Token
	paramsLimit := p.Limit
	paramsOffset := p.Offset

	//トークンからユーザーidを取得する為に利用
	token, err := t.GetByToken(paramsToken)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if token == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Your Token Is Invalid",
			})
	}

	// limitが20になっているかをvalidation
	if paramsLimit != int64(20) {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// offset 負の数になったときにBadRequestを返す
	if paramsOffset < 0 {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	//いいねをすでに送っている人を取得
	userIDs, err := l.FindLikeOnley(token.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	//男か女か情報が必要なのでユーザーの情報を取得する
	user, err := u.GetByUserID(token.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if user == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	//ユーザーが女だった時に検索するユーザーを男とする
	gender := user.GetOppositeGender()

	//明示的に型宣言
	var f entities.Users

	//探す処理
	f, err = u.FindWithCondition(int(paramsLimit), int(paramsOffset), gender, userIDs)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if f == nil {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	sEnt := f.Build()
	return si.NewGetUsersOK().WithPayload(sEnt)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	r := repositories.NewUserRepository()

	// paramsの変数を定義
	paramsToken := p.Token
	paramsUserID := p.UserID

	token, err := t.GetByToken(paramsToken)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if token == nil {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// Bad Requestをどのタイミングで使うかわからないのであとで調査
	// User情報を取得する
	user, err := r.GetByUserID(paramsUserID)
	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if user == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found",
			})
	}

	sEnt := user.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	u := repositories.NewUserRepository()

	// paramsの変数を定義
	paramsToken := p.Params.Token
	paramsUserID := p.UserID

	// ユーザーID取得用
	token, err := t.GetByToken(paramsToken)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if token == nil {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// ユーザーの情報の取得
	us, err := u.GetByUserID(paramsUserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if us == nil {
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	// ユーザーの情報を受け取ったparamsに書き換える
	ProfileUpdate(p.Params, us)

	// tokenのUserIDと受け取ったUserIDが一致しているか確認
	// 一致していなかったら403を返す
	if paramsUserID == token.UserID {
		// 書き換えた情報でデータベースを更新
		err = u.Update(us)
		if err != nil {
			return si.NewPutProfileInternalServerError().WithPayload(
				&si.PutProfileInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}
	} else {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code:    "403",
				Message: "Forbidden",
			})
	}

	// 更新後のユーザー情報の取得
	user, err := u.GetByUserID(paramsUserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if user == nil {
		return si.NewPutProfileBadRequest().WithPayload(
			&si.PutProfileBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	sEnt := user.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}

func ProfileUpdate(p si.PutProfileBody, u *entities.User) {
	// annual income
	u.AnnualIncome = p.AnnualIncome
	// body build
	u.BodyBuild = p.BodyBuild
	// child
	u.Child = p.Child
	// cost of date
	u.CostOfDate = p.CostOfDate
	// drinking
	u.Drinking = p.Drinking
	// education
	u.Education = p.Education
	// height
	u.Height = p.Height
	// holiday
	u.Holiday = p.Holiday
	// home state
	u.HomeState = p.HomeState
	// housework
	u.Housework = p.Housework
	// how to meet
	u.HowToMeet = p.HowToMeet
	// image uri
	u.ImageURI = p.ImageURI
	// introduction
	u.Introduction = p.Introduction
	// job
	u.Job = p.Job
	// marital status
	u.MaritalStatus = p.MaritalStatus
	// nickname
	u.Nickname = p.Nickname
	// nth child
	u.NthChild = p.NthChild
	// residence state
	u.ResidenceState = p.ResidenceState
	// smoking
	u.Smoking = p.Smoking
	// tweet
	u.Tweet = p.Tweet
	// want child
	u.WantChild = p.WantChild
	// when marry
	u.WhenMarry = p.WhenMarry
}
