package user

import (
	"github.com/go-openapi/runtime/middleware"

	"fmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

// 探すAPI
func GetUsers(p si.GetUsersParams) middleware.Responder {
	tokenRepo := repositories.NewUserTokenRepository()
	userRepo := repositories.NewUserRepository()
	likeRepo := repositories.NewUserLikeRepository()

	// tokenが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if tokenOwner == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// トークンの持ち主のIdを元にユーザープロフィールを取得します
	user, err := userRepo.GetByUserID(tokenOwner.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	ids, err := likeRepo.FindLikeAll(tokenOwner.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	ents, err := userRepo.FindWithCondition(int(p.Limit), int(p.Offset), user.GetOppositeGender(), ids)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})

	}

	ent := entities.Users(ents)
	sEnt := ent.Build()
	return si.NewGetUsersOK().WithPayload(sEnt)
}

// ユーザー詳細API
func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	tokenRepo := repositories.NewUserTokenRepository()
	userRepo := repositories.NewUserRepository()

	// tokenが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(p.Token)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if tokenOwner == nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// pathのuserIdを使ってユーザープロフィールを取得します。
	ent, err := userRepo.GetByUserID(p.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if ent == nil {
		return si.NewGetTokenByUserIDNotFound().WithPayload(
			&si.GetTokenByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Not Found",
			})
	}

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

// ユーザー情報更新API
func PutProfile(p si.PutProfileParams) middleware.Responder {
	tokenRepo := repositories.NewUserTokenRepository()
	userRepo := repositories.NewUserRepository()

	// tokenが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(p.Params.Token)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if tokenOwner == nil {
		return si.NewPutProfileUnauthorized().WithPayload(
			&si.PutProfileUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	// bodyの中の値でユーザー情報を更新します。
	fmt.Println(p.Params)
	updateUser := entities.User{
		ID:             p.UserID,
		AnnualIncome:   p.Params.AnnualIncome,
		BodyBuild:      p.Params.BodyBuild,
		Child:          p.Params.Child,
		CostOfDate:     p.Params.CostOfDate,
		Drinking:       p.Params.Drinking,
		Education:      p.Params.Education,
		Height:         p.Params.Height,
		Holiday:        p.Params.Holiday,
		HomeState:      p.Params.HomeState,
		Housework:      p.Params.Housework,
		HowToMeet:      p.Params.HowToMeet,
		ImageURI:       p.Params.ImageURI,
		Introduction:   p.Params.Introduction,
		Job:            p.Params.Introduction,
		MaritalStatus:  p.Params.MaritalStatus,
		Nickname:       p.Params.Nickname,
		NthChild:       p.Params.NthChild,
		ResidenceState: p.Params.ResidenceState,
		Smoking:        p.Params.Smoking,
		Tweet:          p.Params.Tweet,
		WantChild:      p.Params.WantChild,
		WhenMarry:      p.Params.WhenMarry,
	}

	res := userRepo.Update(&updateUser)
	if res != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// pathのuserIdを使って更新後のユーザープロフィールを取得します。
	ent, err := userRepo.GetByUserID(p.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	sEnt := ent.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}
