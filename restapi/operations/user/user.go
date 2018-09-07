package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

var (
	tokenRepo = repositories.NewUserTokenRepository()
	userRepo  = repositories.NewUserRepository()
	likeRepo  = repositories.NewUserLikeRepository()
	imageRepo = repositories.NewUserImageRepository()
)

// 探すAPI
func GetUsers(p si.GetUsersParams) middleware.Responder {
	// 入力値のValidation処理をします。
	limit := int(p.Limit)
	if limit <= 0 {
		return getUsersBadRequestResponses()
	}

	offset := int(p.Offset)
	if offset < 0 {
		return getUsersBadRequestResponses()
	}

	token := p.Token

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return getUsersInternalServerErrorResponse()
	}
	if tokenOwner == nil {
		return getUsersUnauthorizedResponse()
	}

	id := tokenOwner.UserID

	// ユーザーのプロフィールを取得します。
	user, err := userRepo.GetByUserID(id)
	if err != nil {
		return getUsersInternalServerErrorResponse()
	}
	if user == nil {
		return getUsersBadRequestResponses()
	}

	// ユーザーが既にLikeしている or されている状態のUserIDを取得します。
	ids, err := likeRepo.FindLikeAll(id)
	if err != nil {
		return getUsersInternalServerErrorResponse()
	}

	// いいね！を送っていないユーザーを取得します。
	users, err := userRepo.FindWithCondition(limit, offset, user.GetOppositeGender(), ids)
	if err != nil {
		return getUsersInternalServerErrorResponse()
	}

	// いいね！を送っていないユーザーのIDを取得します。
	ids = nil
	for _, user := range users {
		ids = append(ids, user.ID)
	}

	// ユーザーのプロフィール画像を取得します。
	images, err := imageRepo.GetByUserIDs(ids)
	if err != nil {
		return getUsersInternalServerErrorResponse()
	}

	// 取得したお相手のプロフィールとプロフィール画像をマッピングします。
	var ents entities.Users
	for _, user := range users {
		ent := entities.User{}
		for _, image := range images {
			if user.ID == image.UserID {
				ent = user
				ent.ImageURI = image.Path
			}
		}

		ents = append(ents, ent)
	}

	sEnt := ents.Build()
	return si.NewGetUsersOK().WithPayload(sEnt)
}

// ユーザー詳細API
func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	// 入力値のValidation処理をします。
	id := p.UserID
	if id <= 0 {
		return getProfileByUserIDLessThan0BadRequestResponses()
	}

	token := p.Token

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return getProfileByUserIDInternalServerErrorResponse()
	}
	if tokenOwner == nil {
		return getProfileByUserIDUnauthorizedResponse()
	}

	// ユーザーのプロフィール画像を取得します。
	imagePath, err := imageRepo.GetByUserID(id)
	if err != nil {
		return getProfileByUserIDInternalServerErrorResponse()
	}

	// ユーザーのプロフィールを取得します。
	ent, err := userRepo.GetByUserID(id)
	if err != nil {
		return getProfileByUserIDInternalServerErrorResponse()
	}
	if ent == nil {
		return getProfileByUserIDNotFoundResponse()
	}

	ent.ImageURI = imagePath.Path

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

// ユーザー情報更新API
func PutProfile(p si.PutProfileParams) middleware.Responder {
	// 入力値のValidation処理をします。
	id := p.UserID
	if id <= 0 {
		return putProfileLessThan0BadRequestResponse()
	}

	token := p.Params.Token

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return putProfileInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return putProfileUnauthorizedResponse()
	}

	// 他人のデータを改ざんできないようにします。
	if tokenOwner.UserID != id {
		return putProfileForbiddenResponse()
	}

	// ユーザーのプロフィールを取得します。
	user, err := userRepo.GetByUserID(id)
	if err != nil {
		return putProfileInternalServerErrorResponse()
	}

	// 変数をマッピングします。
	user.MappingParams(p.Params)

	// ユーザーのプロフィールを更新します。
	err = userRepo.Update(user)
	if err != nil {
		return putProfileInternalServerErrorResponse()
	}

	// ユーザーのプロフィール画像を取得します。
	imagePath, err := imageRepo.GetByUserID(id)
	if err != nil {
		return getProfileByUserIDInternalServerErrorResponse()
	}

	// 更新後のユーザーのプロフィールを取得します。
	ent, err := userRepo.GetByUserID(id)
	if err != nil {
		return putProfileInternalServerErrorResponse()
	}
	ent.ImageURI = imagePath.Path

	sEnt := ent.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}
