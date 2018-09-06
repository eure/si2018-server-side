package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func getMatchesThrowInternalServerError(fun string, err error) *si.GetMatchesInternalServerError {
	return si.NewGetMatchesInternalServerError().WithPayload(
		&si.GetMatchesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func getMatchesThrowUnauthorized(mes string) *si.GetMatchesUnauthorized {
	return si.NewGetMatchesUnauthorized().WithPayload(
		&si.GetMatchesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func getMatchesThrowBadRequest(mes string) *si.GetMatchesBadRequest {
	return si.NewGetMatchesBadRequest().WithPayload(
		&si.GetMatchesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	var err error
	// トークン認証
	var id int64
	{
		tokenRepo := repositories.NewUserTokenRepository()
		token, err := tokenRepo.GetByToken(p.Token)
		if err != nil {
			return getMatchesThrowInternalServerError("GetByToken", err)
		}
		if token == nil {
			return getMatchesThrowUnauthorized("GetByToken failed")
		}
		id = token.UserID
	}
	// マッチング情報の取得
	var matched []entities.UserMatch
	{
		matchRepo := repositories.NewUserMatchRepository()
		matched, err = matchRepo.FindByUserIDWithLimitOffset(id, int(p.Limit), int(p.Offset))
		if err != nil {
			return getMatchesThrowInternalServerError("FindByUserIDWithLimitOffset", err)
		}
	}
	// 相手の ID の取得
	ids := make([]int64, 0)
	for _, m := range matched {
		ids = append(ids, m.GetPartnerID(id))
	}
	// 相手のユーザー情報を取得
	var users []entities.User
	{
		userRepo := repositories.NewUserRepository()
		users, err = userRepo.FindByIDs(ids)
		if err != nil {
			return getMatchesThrowInternalServerError("FindByIDs", err)
		}
	}
	// 相手の写真を取得
	var images []entities.UserImage
	{
		imageRepo := repositories.NewUserImageRepository()
		images, err = imageRepo.GetByUserIDs(ids)
		if err != nil {
			return getMatchesThrowInternalServerError("GetByUserIDs", err)
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
