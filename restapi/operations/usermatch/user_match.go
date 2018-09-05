package usermatch

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

// マッチング一覧API
func GetMatches(p si.GetMatchesParams) middleware.Responder {
	// 入力値のValidation処理をします。
	limit := int(p.Limit)
	if limit <= 0 {
		return getMatchesBadRequestResponses()
	}

	offset := int(p.Offset)
	if offset < 0 {
		return getMatchesBadRequestResponses()
	}

	token := p.Token

	tokenRepo := repositories.NewUserTokenRepository()
	matchRepo := repositories.NewUserMatchRepository()
	userRepo := repositories.NewUserRepository()

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		return getMatchesInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return getMatchesUnauthorizedResponse()
	}

	id := tokenOwner.UserID

	// ユーザーがマッチングしているお相手の一覧を取得します。
	matches, err := matchRepo.FindByUserIDWithLimitOffset(id, limit, offset)
	if err != nil {
		return getMatchesInternalServerErrorResponse()
	}

	// ユーザーがマッチングしているお相手のIDの配列を取得します。
	var ids []int64
	for _, match := range matches {
		ids = append(ids, match.PartnerID)
	}

	// idの配列からお相手のプロフィールを取得します。
	users, err := userRepo.FindByIDs(ids)
	if err != nil {
		return getMatchesInternalServerErrorResponse()
	}

	var ents entities.MatchUserResponses

	// 取得したお相手のIDとマッチングしているお相手のIDを比較し、マッピングします。
	for _, match := range matches {
		ent := entities.MatchUserResponse{}
		ent.MatchedAt = match.CreatedAt
		for _, user := range users {
			if match.PartnerID == user.ID {
				ent.ApplyUser(user)
			}
		}

		ents = append(ents, ent)
	}

	sEnt := ents.Build()
	return si.NewGetMatchesOK().WithPayload(sEnt)
}

/*			以下　Validationに用いる関数			*/

//	マッチング一覧API
// 	GET {hostname}/api/1.0/matches
func getMatchesInternalServerErrorResponse() middleware.Responder {
	return si.NewGetMatchesInternalServerError().WithPayload(
		&si.GetMatchesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error",
		})
}

func getMatchesUnauthorizedResponse() middleware.Responder {
	return si.NewGetMatchesUnauthorized().WithPayload(
		&si.GetMatchesUnauthorizedBody{
			Code:    "401",
			Message: "Your Token Is Invalid",
		})
}

func getMatchesBadRequestResponses() middleware.Responder {
	return si.NewGetMatchesBadRequest().WithPayload(
		&si.GetMatchesBadRequestBody{
			Code:    "400",
			Message: "Bad Request",
		})
}
