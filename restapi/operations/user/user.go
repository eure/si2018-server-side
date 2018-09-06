package user

import (
	"encoding/json"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	repUser := repositories.NewUserRepository()
	repUserToken := repositories.NewUserTokenRepository()
	repUserLike := repositories.NewUserLikeRepository()

	// TODO: ユーザの一覧だから必要最低限の情報だけをかえしたい
	// 出身地でソートしたい

	// tokenのバリデーション
	err := repUserToken.ValidateToken(p.Token)
	if err != nil {
		return si.NewGetUsersUnauthorized().WithPayload(
			&si.GetUsersUnauthorizedBody{
				Code: "401",
				Message: "Your Token Is Invalid",
					})
	}

	// bad Request
	if p.Limit == 0 {
		return si.NewGetMessagesOK().WithPayload(nil)
	} else if (p.Limit < 1) || (p.Offset < 0) {
		return si.NewGetUsersBadRequest().WithPayload(
			&si.GetUsersBadRequestBody{
				Code: "400",
				Message: "Bad Request",
			})
	}

	// ログインユーザーと反対の性別を取得する
	userToken, _ := repUserToken.GetByToken(p.Token)
	loginUser, _ := repUser.GetByUserID(userToken.UserID)
	oppositeGender := loginUser.GetOppositeGender()

	// ログインユーザーがいいねした人/ログインユーザーをいいねした人のIDを取得
	exceptIds, err := repUserLike.FindLikeAll(userToken.UserID)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// ログインユーザーがいいねしていない人&ログインユーザーをいいねしていない人を取得
	user, err := repUser.FindWithCondition(int(p.Limit), int(p.Offset), oppositeGender, exceptIds)
	if err != nil {
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// Userの配列をUsersにキャストする
	users := entities.Users(user)

	usersModel := users.Build()

	return si.NewGetUsersOK().WithPayload(usersModel)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	repUser := repositories.NewUserRepository()
	repUserToken := repositories.NewUserTokenRepository()

	// tokenのバリデーション
	err := repUserToken.ValidateToken(p.Token)
	if err != nil {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code: "401",
				Message: "Your Token Is Invalid",
			})
	}

	// ユーザの取得
	user, err := repUser.GetByUserID(p.UserID)
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
				Code:    "400",
				Message: "User Not Found",
			})
	}

	userModel := user.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&userModel)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	repUser := repositories.NewUserRepository()
	repUserToken := repositories.NewUserTokenRepository()

	// tokenのバリデーション
	err := repUserToken.ValidateToken(p.Params.Token)
	if err != nil {
		return si.NewGetProfileByUserIDUnauthorized().WithPayload(
			&si.GetProfileByUserIDUnauthorizedBody{
				Code: "401",
				Message: "Your Token Is Invalid",
			})
	}

	// Forbidden
	userToken, _ := repUserToken.GetByToken(p.Params.Token)
	if userToken.UserID != p.UserID {
		return si.NewPutProfileForbidden().WithPayload(
			&si.PutProfileForbiddenBody{
				Code: "403",
				Message: "Forbidden",
			})
	}

	// 更新用のUserを作成
	var updateUser entities.User
	updateUser.ID = p.UserID

	// パラメーターの値をupdateUserに入れる
	bindParams(p.Params, &updateUser)

	// Userを更新
	err = repUser.Update(&updateUser)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// 更新したユーザーを再取得
	updatedUser, err := repUser.GetByUserID(userToken.UserID)
	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	updateUserModel := updatedUser.Build()

	return si.NewPutProfileOK().WithPayload(&updateUserModel)
}

// private
func bindParams(p si.PutProfileBody, entUser *entities.User ){
	// paramsをjsonに出力
	params, _ := p.MarshalBinary()
	// userEntにjson変換したparamを入れる
	json.Unmarshal(params, &entUser)

	entUser.HowToMeet = p.HowToMeet
	entUser.AnnualIncome = p.AnnualIncome
	entUser.CostOfDate = p.CostOfDate
	entUser.NthChild = p.NthChild
	entUser.ResidenceState = p.ResidenceState
}
