package usermatch

import (
	"fmt"
	"log"

	"sync"

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

	var wg sync.WaitGroup
	usersChan := make(chan map[int64]entities.User)
	imagesChan := make(chan map[int64]entities.UserImage)
	errChan := make(chan error)

	go func() {
		defer wg.Done()
		fmt.Println("Start Get Profile")
		usrs, err := userRepo.FindByIDs(ids)
		if err != nil {
			log.Fatal(err)
			errChan <- err
		}

		mpU := map[int64]entities.User{}
		for _, u := range usrs {
			if u.ID > 0 {
				mpU[u.ID] = u
			}
		}

		fmt.Println("End Get Profile")
		fmt.Println(mpU)
		usersChan <- mpU
		close(usersChan)
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Start Get ImagePath")
		imgs, err := imageRepo.GetByUserIDs(ids)
		if err != nil {
			log.Fatal(err)
			errChan <- err
		}

		mpI := map[int64]entities.UserImage{}
		for _, img := range imgs {
			if img.UserID > 0 {
				mpI[img.UserID] = img
			}
		}

		fmt.Println("End Get ImagePath")
		fmt.Println(mpI)
		imagesChan <- mpI
		close(imagesChan)
	}()

	fmt.Println("END")
	err = <-errChan
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
