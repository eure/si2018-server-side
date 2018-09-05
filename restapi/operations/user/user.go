package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
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

	tokenRepo := repositories.NewUserTokenRepository()
	userRepo := repositories.NewUserRepository()
	likeRepo := repositories.NewUserLikeRepository()

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return getUsersInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return getUsersUnauthorizedResponse()
	}

	id := tokenOwner.UserID

	// ユーザーのプロフィールを取得します。
	user, err := userRepo.GetByUserID(id)
	if err != nil {
		return getUsersInternalServerErrorResponse()
	} else if user == nil {
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

	ent := entities.Users(users)
	sEnt := ent.Build()
	return si.NewGetUsersOK().WithPayload(sEnt)
}

// ユーザー詳細API
func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	// 入力値のValidation処理をします。
	id := p.UserID
	if id <= 0 {
		return getProfileByUserIDBadRequestResponses()
	}

	token := p.Token

	tokenRepo := repositories.NewUserTokenRepository()
	userRepo := repositories.NewUserRepository()

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return getProfileByUserIDInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return getProfileByUserIDUnauthorizedResponse()
	}

	// ユーザーのプロフィールを取得します。
	ent, err := userRepo.GetByUserID(id)
	if err != nil {
		return getProfileByUserIDInternalServerErrorResponse()
	} else if ent == nil {
		return getProfileByUserIDNotFoundResponse()
	}

	sEnt := ent.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

// ユーザー情報更新API
func PutProfile(p si.PutProfileParams) middleware.Responder {
	// 入力値のValidation処理をします。
	id := p.UserID
	if id <= 0 {
		return putProfileBadRequestResponse()
	}

	token := p.Params.Token

	tokenRepo := repositories.NewUserTokenRepository()
	userRepo := repositories.NewUserRepository()

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

	// 更新後のユーザーのプロフィールを取得します。
	ent, err := userRepo.GetByUserID(p.UserID)
	if err != nil {
		return putProfileInternalServerErrorResponse()
	}

	sEnt := ent.Build()
	return si.NewPutProfileOK().WithPayload(&sEnt)
}

/*			以下　Validationに用いる関数			*/

//	探すAPI
// 	GET {hostname}/api/1.0/users
func getUsersInternalServerErrorResponse() middleware.Responder {
	return si.NewGetUsersInternalServerError().WithPayload(
		&si.GetUsersInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getUsersUnauthorizedResponse() middleware.Responder {
	return si.NewGetUsersUnauthorized().WithPayload(
		&si.GetUsersUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func getUsersBadRequestResponses() middleware.Responder {
	return si.NewGetUsersBadRequest().WithPayload(
		&si.GetUsersBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

//	ユーザー詳細API
//	GET {hostname}/api/1.0/users/{userID}
func getProfileByUserIDInternalServerErrorResponse() middleware.Responder {
	return si.NewGetProfileByUserIDInternalServerError().WithPayload(
		&si.GetProfileByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getProfileByUserIDNotFoundResponse() middleware.Responder {
	return si.NewGetProfileByUserIDNotFound().WithPayload(
		&si.GetProfileByUserIDNotFoundBody{
			Code:    "404",
			Message: "User Not Found",
		})
}

func getProfileByUserIDUnauthorizedResponse() middleware.Responder {
	return si.NewGetProfileByUserIDUnauthorized().WithPayload(
		&si.GetProfileByUserIDUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func getProfileByUserIDBadRequestResponses() middleware.Responder {
	return si.NewGetProfileByUserIDBadRequest().WithPayload(
		&si.GetProfileByUserIDBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}

//	ユーザー情報更新API
//	PUT {hostname}/api/1.0/users/{userID}
func putProfileInternalServerErrorResponse() middleware.Responder {
	return si.NewPutProfileInternalServerError().WithPayload(
		&si.PutProfileInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func putProfileForbiddenResponse() middleware.Responder {
	return si.NewPutProfileForbidden().WithPayload(
		&si.PutProfileForbiddenBody{
			Code:    "403",
			Message: "Forbidden",
		})
}

func putProfileUnauthorizedResponse() middleware.Responder {
	return si.NewPutProfileUnauthorized().WithPayload(
		&si.PutProfileUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func putProfileBadRequestResponse() middleware.Responder {
	return si.NewPutProfileBadRequest().WithPayload(
		&si.PutProfileBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}
