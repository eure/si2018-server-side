package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func getUsersThrowInternalServerError(fun string, err error) *si.GetUsersInternalServerError {
	return si.NewGetUsersInternalServerError().WithPayload(
		&si.GetUsersInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func getUsersThrowUnauthorized(mes string) *si.GetUsersUnauthorized {
	return si.NewGetUsersUnauthorized().WithPayload(
		&si.GetUsersUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func getUsersThrowBadRequest(mes string) *si.GetUsersBadRequest {
	return si.NewGetUsersBadRequest().WithPayload(
		&si.GetUsersBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetUsers(p si.GetUsersParams) middleware.Responder {
	var err error
	userRepo := repositories.NewUserRepository()

	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Token)
		if err != nil {
			return getUsersThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return getUsersThrowUnauthorized("GetByToken failed")
		}
		id = token.UserID
	}
	// 異性のみを検索するために, 性別情報が必要
	var oppositeGender string
	{
		user, err := userRepo.GetByUserID(id)
		if err != nil {
			return getUsersThrowInternalServerError("GetByUserID", err)
		}
		if user == nil {
			return getUsersThrowBadRequest("GetByUserID failed")
		}
		oppositeGender = user.GetOppositeGender()
	}
	// いいねまたはマッチ済みの相手を取得
	ids := make([]int64, 0)
	{
		likeRepo := repositories.NewUserLikeRepository()
		matchRepo := repositories.NewUserMatchRepository()
		idmap := make(map[int64]bool)
		like, err := likeRepo.FindLikeAll(id)
		if err != nil {
			return getUsersThrowInternalServerError("FindLikeAll", err)
		}
		for _, id := range like {
			idmap[id] = true
		}
		matched, err := matchRepo.FindAllByUserID(id)
		if err != nil {
			return getUsersThrowInternalServerError("FindAllByUserID", err)
		}
		for _, id := range matched {
			idmap[id] = true
		}
		for k := range idmap {
			ids = append(ids, k)
		}
	}
	// 相手を探す
	partnerInfos, err := userRepo.FindWithCondition(int(p.Limit), int(p.Offset), oppositeGender, ids)
	if err != nil {
		return getUsersThrowInternalServerError("FindWithCondition", err)
	}
	if partnerInfos == nil {
		return getUsersThrowBadRequest("FindWithCondition")
	}
	// 相手の画像を取得する
	var partnerImages []entities.UserImage
	{
		imageRepo := repositories.NewUserImageRepository()
		partnerIDs := make([]int64, 0)
		for _, u := range partnerInfos {
			partnerIDs = append(partnerIDs, u.ID)
		}
		partnerImages, err = imageRepo.GetByUserIDs(partnerIDs)
		if err != nil {
			return getUsersThrowInternalServerError("GetByUserIDs", err)
		}
		if len(partnerImages) != len(partnerInfos) {
			return getUsersThrowBadRequest("GetByUserIDs failed")
		}
	}
	for i := range partnerInfos {
		partnerInfos[i].ImageURI = partnerImages[i].Path
	}
	partners := entities.Users(partnerInfos)

	return si.NewGetUsersOK().WithPayload(partners.Build())
}

func getProfileByUserIDThrowInternalServerError(fun string, err error) *si.GetProfileByUserIDInternalServerError {
	return si.NewGetProfileByUserIDInternalServerError().WithPayload(
		&si.GetProfileByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func getProfileByUserIDThrowUnauthorized(mes string) *si.GetProfileByUserIDUnauthorized {
	return si.NewGetProfileByUserIDUnauthorized().WithPayload(
		&si.GetProfileByUserIDUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func getProfileByUserIDThrowBadRequest(mes string) *si.GetProfileByUserIDBadRequest {
	return si.NewGetProfileByUserIDBadRequest().WithPayload(
		&si.GetProfileByUserIDBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func getProfileByUserIDThrowNotFound(mes string) *si.GetProfileByUserIDNotFound {
	return si.NewGetProfileByUserIDNotFound().WithPayload(
		&si.GetProfileByUserIDNotFoundBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	userRepo := repositories.NewUserRepository()

	// トークン認証
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Token)
		if err != nil {
			return getProfileByUserIDThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return getProfileByUserIDThrowUnauthorized("GetByToken failed")
		}
	}
	user, err := userRepo.GetByUserID(p.UserID)
	if err != nil {
		return getProfileByUserIDThrowInternalServerError("GetByUserID", err)
	}
	if user == nil {
		return getProfileByUserIDThrowBadRequest("GetByUserID failed")
	}
	// 画像を取得する
	var image *entities.UserImage
	{
		imageRepo := repositories.NewUserImageRepository()
		image, err = imageRepo.GetByUserID(p.UserID)
		if err != nil {
			return getUsersThrowInternalServerError("GetByUserID", err)
		}
		if image == nil {
			return getUsersThrowBadRequest("GetByUserID failed")
		}
	}
	user.ImageURI = image.Path

	sEnt := user.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func putProfileThrowInternalServerError(fun string, err error) *si.PutProfileInternalServerError {
	return si.NewPutProfileInternalServerError().WithPayload(
		&si.PutProfileInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func putProfileThrowUnauthorized(mes string) *si.PutProfileUnauthorized {
	return si.NewPutProfileUnauthorized().WithPayload(
		&si.PutProfileUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func putProfileThrowBadRequest(mes string) *si.PutProfileBadRequest {
	return si.NewPutProfileBadRequest().WithPayload(
		&si.PutProfileBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func putProfileThrowForbidden(mes string) *si.PutProfileForbidden {
	return si.NewPutProfileForbidden().WithPayload(
		&si.PutProfileForbiddenBody{
			Code:    "403",
			Message: "Forbidden. (他の人のプロフィールは更新できません.): " + mes,
		})
}

func ApplyParams(user *entities.User, params si.PutProfileBody) {
	if user.MaritalStatus != "独身(未婚)" && params.MaritalStatus == "独身(未婚)" {
		// 未婚でない人が未婚になることはありえないはず
		// システム上、結婚歴で嘘をつくことを認めるかは議論の余地あり
		// とりあえず今回は見逃す
	}
	user.AnnualIncome = params.AnnualIncome
	user.BodyBuild = params.BodyBuild
	user.Child = params.Child
	user.CostOfDate = params.CostOfDate
	user.Drinking = params.Drinking
	user.Education = params.Education
	user.Height = params.Height
	user.Holiday = params.Holiday
	user.HomeState = params.HomeState
	user.Housework = params.Housework
	user.HowToMeet = params.HowToMeet
	// 画像の更新は post images で行う
	// params から除外すべき
	// u.ImageURI = params.ImageURI
	user.Introduction = params.Introduction
	user.Job = params.Job
	user.MaritalStatus = params.MaritalStatus
	user.Nickname = params.Nickname
	user.NthChild = params.NthChild
	user.ResidenceState = params.ResidenceState
	user.Smoking = params.Smoking
	user.Tweet = params.Tweet
	user.WantChild = params.WantChild
	user.WhenMarry = params.WhenMarry
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	var err error
	userRepo := repositories.NewUserRepository()

	// トークン認証
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Params.Token)
		if err != nil {
			return putProfileThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return putProfileThrowUnauthorized("GetByToken failed")
		}
		if token.UserID != p.UserID {
			return putProfileThrowForbidden("Token does not match")
		}
	}
	// ユーザー情報を取得して更新を反映させる
	var user *entities.User
	{
		user, err = userRepo.GetByUserID(p.UserID)
		if err != nil {
			return putProfileThrowInternalServerError("GetByUserID", err)
		}
		if user == nil {
			return putProfileThrowBadRequest("GetByUserID failed")
		}
		ApplyParams(user, p.Params)
	}
	err = userRepo.Update(user)
	if err != nil {
		return putProfileThrowInternalServerError("Update", err)
	}
	// 更新後のユーザーを取得し直す (これをしないと, p.Params に nil があるときに整合しない)
	user, err = userRepo.GetByUserID(p.UserID)
	if err != nil {
		return putProfileThrowInternalServerError("GetByUserID", err)
	}
	if user == nil {
		return putProfileThrowBadRequest("GetByUserID failed")
	}
	// 画像を取得する
	var image *entities.UserImage
	{
		imageRepo := repositories.NewUserImageRepository()
		image, err = imageRepo.GetByUserID(p.UserID)
		if err != nil {
			return putProfileThrowInternalServerError("GetByUserID", err)
		}
		if image == nil {
			return putProfileThrowBadRequest("GetByUserID failed")
		}
	}
	user.ImageURI = image.Path
	sEnt := user.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}
