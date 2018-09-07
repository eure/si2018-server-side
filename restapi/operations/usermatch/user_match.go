package usermatch

import (
	"log"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

var (
	tokenRepo = repositories.NewUserTokenRepository()
	matchRepo = repositories.NewUserMatchRepository()
	userRepo  = repositories.NewUserRepository()
	imageRepo = repositories.NewUserImageRepository()
)

// マッチング一覧API
func GetMatches(p si.GetMatchesParams) middleware.Responder {
	// 入力値のValidation処理をします。
	limit := int(p.Limit)
	if limit < 1 {
		return getMatchesLimitBadRequestResponses()
	}

	offset := int(p.Offset)
	if offset < 0 {
		return getMatchesOffsetBadRequestResponses()
	}

	token := p.Token

	// トークンが有効であるか検証します。
	tokenOwner, err := tokenRepo.GetByToken(token)
	if err != nil {
		log.Fatal(err)
		return getMatchesInternalServerErrorResponse()
	} else if tokenOwner == nil {
		return getMatchesUnauthorizedResponse()
	}

	id := tokenOwner.UserID

	// ユーザーがマッチングしているお相手の一覧を取得します。
	matches, err := matchRepo.FindByUserIDWithLimitOffset(id, limit, offset)
	if err != nil {
		log.Fatal(err)
		return getMatchesInternalServerErrorResponse()
	}

	// matchesの長さのスライスを作成し、ユーザーがマッチングしているお相手のIDの配列を取得します。
	ids := make([]int64, len(matches))
	for _, match := range matches {
		if match.UserID == id {
			ids = append(ids, match.PartnerID)
		}
		if match.PartnerID == id {
			ids = append(ids, match.UserID)
		}
	}

	usersChan := make(chan map[int64]entities.User, 1)
	imagesChan := make(chan map[int64]entities.UserImage, 1)
	errChan1 := make(chan error, 1)
	errChan2 := make(chan error, 1)

	go func(result chan error) {
		usrs, err := userRepo.FindByIDs(ids)
		if err != nil {
			log.Fatal(err)
			result <- err
		}

		mpU := map[int64]entities.User{}
		for _, u := range usrs {
			if u.ID > 0 {
				mpU[u.ID] = u
			}
		}

		usersChan <- mpU
		result <- nil

	}(errChan1)

	go func(result chan error) {
		imgs, err := imageRepo.GetByUserIDs(ids)
		if err != nil {
			log.Fatal(err)
			result <- err
		}

		mpI := map[int64]entities.UserImage{}
		for _, img := range imgs {
			if img.UserID > 0 {
				mpI[img.UserID] = img
			}
		}

		imagesChan <- mpI
		result <- nil

	}(errChan2)

	err = <-errChan1
	if err != nil {
		return getMatchesInternalServerErrorResponse()
	}

	err = <-errChan2
	if err != nil {
		return getMatchesInternalServerErrorResponse()
	}

	users := <-usersChan
	images := <-imagesChan

	var ents entities.MatchUserResponses
	for _, match := range matches {
		ent := entities.MatchUserResponse{}
		if match.UserID == id {
			ent.ApplyUser(users[match.PartnerID])
			ent.ImageURI = images[match.PartnerID].Path
			ent.MatchedAt = match.CreatedAt
		}
		if match.PartnerID == id {
			ent.ApplyUser(users[match.UserID])
			ent.ImageURI = images[match.UserID].Path
			ent.MatchedAt = match.CreatedAt
		}
		ents = append(ents, ent)
	}

	sEnt := ents.Build()
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
