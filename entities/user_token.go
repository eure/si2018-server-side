package entities

import (
	"github.com/eure/si2018-server-side/models"
	"github.com/go-openapi/strfmt"
)

type UserToken struct {
	UserID    int64           `xorm:"user_id"`
	Token     string          `xorm:"token"`
	CreatedAt strfmt.DateTime `xorm:"created_at created"`
	UpdatedAt strfmt.DateTime `xorm:"updated_at updated"`
}

func (u UserToken) Build() models.UserToken {
	return models.UserToken{
		UserID:    u.UserID,
		Token:     u.Token,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
