package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

// DB アクセス: 4 回
// 計算量: O(N)
func GetMatches(p si.GetMatchesParams) middleware.Responder {
	var err error
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Token)
		if err != nil {
			return si.GetMatchesThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return si.GetMatchesThrowUnauthorized("GetByToken failed")
		}
		id = token.UserID
	}
	// マッチング情報の取得
	var matched []entities.UserMatch
	{
		matchRepo := repositories.NewUserMatchRepository()
		matched, err = matchRepo.FindByUserIDWithLimitOffset(id, int(p.Limit), int(p.Offset))
		if err != nil {
			return si.GetMatchesThrowInternalServerError("FindByUserIDWithLimitOffset", err)
		}
	}
	// 相手の ID の取得
	ids := make([]int64, 0)
	for _, m := range matched {
		ids = append(ids, m.GetPartnerID(id))
	}
	mapping := make(map[int64]int)
	for i, m := range matched {
		mapping[m.UserID] = i
	}
	count := len(matched)
	// 相手のユーザー情報を取得
	users := make([]entities.User, count)
	{
		userRepo := repositories.NewUserRepository()
		shuffledUsers, err := userRepo.FindByIDs(ids)
		if err != nil {
			return si.GetMatchesThrowInternalServerError("FindByIDs", err)
		}
		if len(shuffledUsers) != count {
			return si.GetMatchesThrowBadRequest("FindByIDs failed")
		}
		for _, u := range shuffledUsers {
			users[mapping[u.ID]] = u
		}
	}
	// 相手の写真を取得
	images := make([]entities.UserImage, count)
	{
		imageRepo := repositories.NewUserImageRepository()
		shuffledImages, err := imageRepo.GetByUserIDs(ids)
		if err != nil {
			return si.GetMatchesThrowInternalServerError("GetByUserIDs", err)
		}
		if len(shuffledImages) != count {
			return si.GetMatchesThrowBadRequest("GetByUserIDs failed")
		}
		for _, im := range shuffledImages {
			images[mapping[im.UserID]] = im
		}
	}
	// 以上の情報をまとめる
	sEnt := make([]*models.MatchUserResponse, 0)
	for i, m := range matched {
		response := entities.MatchUserResponse{MatchedAt: m.CreatedAt}
		response.ApplyUser(users[i])
		response.ImageURI = images[i].Path
		swaggerMatch := response.Build()
		sEnt = append(sEnt, &swaggerMatch)
	}
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
