package usermatch

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/go-openapi/runtime/middleware"
	"github.com/eure/si2018-server-side/entities"
)

func GetMatches(p si.GetMatchesParams) middleware.Responder {
	t := repositories.NewUserTokenRepository()
	m := repositories.NewUserMatchRepository()
	u := repositories.NewUserRepository()

	ut, _ := t.GetByToken(p.Token)

	// マッチ済みのユーザーを取得
	match, _ := m.FindByUserIDWithLimitOffset(ut.UserID, int(p.Limit), int(p.Offset))

	// ApplyUserの形にしてBuildしないといけないので
	// ApplyUserの型に変換する
	var mr entities.MatchUserResponses
	for _, matchUser := range match{
		// 構造体の初期化(単体)
		r := entities.MatchUserResponse{}
		user, _ := u.GetByUserID(matchUser.PartnerID)
		r.ApplyUser(*user)
		mr = append(mr, r)
	}

  sEnt := mr.Build()
	return si.NewGetMatchesOK().WithPayload(sEnt)
}
