package user

import (
	"encoding/json"
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	repoUser := repositories.NewUserRepository()
	repoUserToken := repositories.NewUserTokenRepository()
	repoUserLike := repositories.NewUserLikeRepository()

	// ログインユーザーと反対の性別を取得する
	entUserToken, _ := repoUserToken.GetByToken(p.Token)
	entUser, _ := repoUser.GetByUserID(entUserToken.UserID)
	opposite_gender := entUser.GetOppositeGender()

	// idsには取得対象に含めないUserIDを入れる (いいね/マッチ/ブロック済みなど) いいねやマッチした人、ブロックした人のidを取ってくる
	//いいねした/された人のidを持ってくる
	except_ids, _ := repoUserLike.FindLikeAll(entUserToken.UserID)

	userEnts, _ := repoUser.FindWithCondition(int(p.Limit), int(p.Offset), opposite_gender, except_ids)

	usersEnt := entities.Users(userEnts)

	users := usersEnt.Build()

	return si.NewGetUsersOK().WithPayload(users)
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	r := repositories.NewUserRepository()

	userEnt, err := r.GetByUserID(p.UserID)

	if err != nil {
		return si.NewGetProfileByUserIDInternalServerError().WithPayload(
			&si.GetProfileByUserIDInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if userEnt == nil {
		return si.NewGetProfileByUserIDNotFound().WithPayload(
			&si.GetProfileByUserIDNotFoundBody{
				Code:    "404",
				Message: "User Profile Not Found",
			})
	}

	user := userEnt.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&user)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	r := repositories.NewUserRepository()

	userEnt := entities.User{ID: p.UserID}

	BindParams(p.Params, &userEnt)

	err := r.Update(&userEnt)

	if err != nil {
		return si.NewPutProfileInternalServerError().WithPayload(
			&si.PutProfileInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	user := userEnt.Build()

	return si.NewPutProfileOK().WithPayload(&user)
}

// private
func BindParams(p si.PutProfileBody, userEnt *entities.User ){
	// paramsをjsonに出力
	params, _ := p.MarshalBinary()
	// userEntにjson変換したparamを入れる
	json.Unmarshal(params, &userEnt)

	userEnt.HowToMeet = p.HowToMeet
	userEnt.AnnualIncome = p.AnnualIncome
	userEnt.CostOfDate = p.CostOfDate
	userEnt.NthChild = p.NthChild
	userEnt.ResidenceState = p.ResidenceState
}
