package user

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
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

	//ユーザーが男だった時に検索するユーザーを女とする
	var g string
	if user.Gender == "M" {
		g = "F"
	}else{
		g = "M"
	}
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
	return si.NewPutProfileOK()
}
