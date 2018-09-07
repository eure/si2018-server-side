package entities

import (
	"time"
	"github.com/eure/si2018-server-side/models"
	"github.com/go-openapi/strfmt"
)


type UserMatch struct {
	UserID    int64           `xorm:"user_id"`
	PartnerID int64           `xorm:"partner_id"`
	CreatedAt strfmt.DateTime `xorm:"created_at"`
	UpdatedAt strfmt.DateTime `xorm:"updated_at"`
}

func NewUserMatch(userID int64, partnerID int64) UserMatch {
	now := strfmt.DateTime(time.Now())
	return UserMatch {
		UserID:    userID,
		PartnerID: partnerID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u UserMatch) Build() models.UserMatch {
	return models.UserMatch{
		UserID:    u.UserID,
		PartnerID: u.PartnerID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type UserMatches []UserMatch

func (uu *UserMatches) Build() []*models.UserMatch {
	var sUsers []*models.UserMatch

	for _, u := range *uu {
		swaggerUser := u.Build()
		sUsers = append(sUsers, &swaggerUser)
	}
	return sUsers
}
