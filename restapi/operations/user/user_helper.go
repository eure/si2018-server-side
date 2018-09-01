package user

import (
	"github.com/eure/si2018-server-side/entities"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func BuildUserEntityByModel(meID int64, p si.PutProfileBody) entities.User {
	return entities.User{
		ID: meID,

		Nickname:       p.Nickname,
		ImageURI:       p.ImageURI,
		Tweet:          p.Tweet,
		Introduction:   p.Introduction,
		ResidenceState: p.ResidenceState,
		HomeState:      p.HomeState,
		Education:      p.Education,
		Job:            p.Job,
		AnnualIncome:   p.AnnualIncome,
		Height:         p.Height,
		BodyBuild:      p.BodyBuild,
		MaritalStatus:  p.MaritalStatus,
		Child:          p.Child,
		WhenMarry:      p.WhenMarry,
		WantChild:      p.WantChild,
		Smoking:        p.Smoking,
		Drinking:       p.Drinking,
		Holiday:        p.Holiday,
		HowToMeet:      p.HowToMeet,
		CostOfDate:     p.CostOfDate,
		NthChild:       p.NthChild,
		Housework:      p.Housework,
	}
}

// プロフィール画像を挿入
func SetUsersImage(users []entities.User) ([]entities.User, error) {
	userList := make([]entities.User, len(users))
	var err error

	return userList, err
}

// プロフィール画像を挿入
func SetUserImage(user entities.User) (entities.User, error) {
	var err error
	return user, err
}
