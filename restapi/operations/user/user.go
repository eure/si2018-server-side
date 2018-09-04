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

	//トークンからユーザーidを取得する為に利用
	ut, _ := t.GetByToken(p.Token)
	//いいねをすでに送っている人を取得
	ids, _ := l.FindLikeAll(ut.UserID)
	//男か女か情報が必要なのでユーザーの情報を取得する
	user, _ := u.GetByUserID(ut.UserID)

	//ユーザーが女だった時に検索するユーザーを男とする
	g := user.GetOppositeGender()

	//明示的に型宣言
	var f entities.Users
	//探す処理
	f, _ = u.FindWithCondition(int(p.Limit), int(p.Offset), g, ids)

	sEnt := f.Build()
	return si.NewGetUsersOK().WithPayload(sEnt)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	// t := repositories.NewUserTokenRepository()
	r := repositories.NewUserRepository()

	// token, _ = t.GetByToken(p.Token)

	ent, err := r.GetByUserID(p.UserID)

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
				Message: "User Token Not Found",
			})
	}

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	u := repositories.NewUserRepository()

	// ユーザーID取得用
	token, _ := t.GetByToken(p.Params.Token)
	// ユーザーの情報の取得
	us, _ := u.GetByUserID(p.UserID)
	// ユーザーの情報を受け取ったparamsに書き換える
	ProfileUpdate(p.Params, us)

	// tokenのUserIDと受け取ったUserIDが一致しているか確認
	// 一致していなかったら403を返す
	if p.UserID == token.UserID {
		// 書き換えた情報でデータベースを更新
		err := u.Update(us)
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
	user, _ := u.GetByUserID(p.UserID)
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
