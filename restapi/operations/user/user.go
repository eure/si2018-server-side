package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

// DB アクセス: 6 回
// 計算量: O(N)
func GetUsers(p si.GetUsersParams) middleware.Responder {
	var err error
	userRepo := repositories.NewUserRepository()

	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Token)
		if err != nil {
			return si.GetUsersThrowInternalServerError(err)
		}
		if token == nil {
			return si.GetUsersThrowUnauthorized()
		}
		id = token.UserID
	}
	// 異性のみを検索するために, 性別情報が必要
	var oppositeGender string
	{
		user, err := userRepo.GetByUserID(id)
		if err != nil {
			return si.GetUsersThrowInternalServerError(err)
		}
		if user == nil {
			return si.GetUsersThrowInternalServerError(nil)
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
			return si.GetUsersThrowInternalServerError(err)
		}
		for _, id := range like {
			idmap[id] = true
		}
		matched, err := matchRepo.FindAllByUserID(id)
		if err != nil {
			return si.GetUsersThrowInternalServerError(err)
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
		return si.GetUsersThrowInternalServerError(err)
	}
	count := len(partnerInfos)
	// あとで順番を調整するのに用いる
	mapping := make(map[int64]int)
	for i, m := range partnerInfos {
		mapping[m.ID] = i
	}
	// 相手の画像を取得する
	partnerImages := make([]entities.UserImage, count)
	{
		imageRepo := repositories.NewUserImageRepository()
		partnerIDs := make([]int64, 0)
		for _, u := range partnerInfos {
			partnerIDs = append(partnerIDs, u.ID)
		}
		shuffledPartnerImages, err := imageRepo.GetByUserIDs(partnerIDs)
		if err != nil {
			return si.GetUsersThrowInternalServerError(err)
		}
		if len(shuffledPartnerImages) != count {
			return si.GetUsersThrowInternalServerError(nil)
		}
		// 正しい順番に直す
		for _, im := range shuffledPartnerImages {
			partnerImages[mapping[im.UserID]] = im
		}
	}
	for i := range partnerInfos {
		partnerInfos[i].ImageURI = partnerImages[i].Path
	}
	partners := entities.Users(partnerInfos)

	return si.NewGetUsersOK().WithPayload(partners.Build())
}

// DB アクセス: 3 回
// 計算量: O(1)
func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	userRepo := repositories.NewUserRepository()

	// トークン認証
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Token)
		if err != nil {
			return si.GetProfileByUserIDThrowInternalServerError(err)
		}
		if token == nil {
			return si.GetProfileByUserIDThrowUnauthorized()
		}
	}
	user, err := userRepo.GetByUserID(p.UserID)
	if err != nil {
		return si.GetProfileByUserIDThrowInternalServerError(err)
	}
	if user == nil {
		return si.GetProfileByUserIDThrowInternalServerError(nil)
	}
	// 画像を取得する
	var image *entities.UserImage
	{
		imageRepo := repositories.NewUserImageRepository()
		image, err = imageRepo.GetByUserID(p.UserID)
		if err != nil {
			return si.GetUsersThrowInternalServerError(err)
		}
		if image == nil {
			return si.GetUsersThrowInternalServerError(nil)
		}
	}
	user.ImageURI = image.Path

	sEnt := user.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
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

// DB アクセス: 5 回
// 計算量: O(1)
func PutProfile(p si.PutProfileParams) middleware.Responder {
	var err error
	userRepo := repositories.NewUserRepository()

	// トークン認証
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Params.Token)
		if err != nil {
			return si.PutProfileThrowInternalServerError(err)
		}
		if token == nil {
			return si.PutProfileThrowUnauthorized()
		}
		if token.UserID != p.UserID {
			return si.PutProfileThrowForbidden()
		}
	}
	// ユーザー情報を取得して更新を反映させる
	var user *entities.User
	{
		user, err = userRepo.GetByUserID(p.UserID)
		if err != nil {
			return si.PutProfileThrowInternalServerError(err)
		}
		if user == nil {
			return si.PutProfileThrowInternalServerError(nil)
		}
		ApplyParams(user, p.Params)
	}
	err = userRepo.Update(user)
	if err != nil {
		return si.PutProfileThrowInternalServerError(err)
	}
	// 更新後のユーザーを取得し直す (これをしないと, p.Params に nil があるときに整合しない)
	user, err = userRepo.GetByUserID(p.UserID)
	if err != nil {
		return si.PutProfileThrowInternalServerError(err)
	}
	if user == nil {
		return si.PutProfileThrowInternalServerError(nil)
	}
	// 画像を取得する
	var image *entities.UserImage
	{
		imageRepo := repositories.NewUserImageRepository()
		image, err = imageRepo.GetByUserID(p.UserID)
		if err != nil {
			return si.PutProfileThrowInternalServerError(err)
		}
		if image == nil {
			return si.PutProfileThrowInternalServerError(nil)
		}
	}
	user.ImageURI = image.Path
	sEnt := user.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}
