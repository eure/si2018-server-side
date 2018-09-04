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

	// ログインユーザーと反対の性別を取得する
	entUserToken, _ := repUserToken.GetByToken(p.Token)
	entUser, _ := repUser.GetByUserID(entUserToken.UserID)
	opposite_gender := entUser.GetOppositeGender()

	// idsには取得対象に含めないUserIDを入れる (いいね/マッチ/ブロック済みなど) いいねやマッチした人、ブロックした人のidを取ってくる
	//いいねした/された人のidを持ってくる
	except_ids, _ := repUserLike.FindLikeAll(entUserToken.UserID)

	entsUser, _ := repUser.FindWithCondition(int(p.Limit), int(p.Offset), opposite_gender, except_ids)

	entUsers := entities.Users(entsUser)

	users := entUsers.Build()

	return si.NewGetUsersOK().WithPayload(users)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	repUser := repositories.NewUserRepository()

	entUser, err := repUser.GetByUserID(p.UserID)

	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if entUser == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Profile Not Found",
			})
	}

	user := entUser.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&user)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	repUser := repositories.NewUserRepository()

	entUser := entities.User{ID: p.UserID}

	BindParams(p.Params, &entUser)

	err := repUser.Update(&entUser)

	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	user := entUser.Build()

	return si.NewPutProfileOK().WithPayload(&user)
}

// private
func BindParams(p si.PutProfileBody, entUser *entities.User ){
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
