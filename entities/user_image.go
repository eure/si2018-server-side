package entities

import (
	"github.com/eure/si2018-server-side/models"
	"github.com/go-openapi/strfmt"
)

type UserImage struct {
	UserID    int64           `xorm:"user_id"`
	Path      string          `xorm:"path"`
	CreatedAt strfmt.DateTime `xorm:"created_at created updated"`
	UpdatedAt strfmt.DateTime `xorm:"updated_at created updated"`
}

func (u UserImage) Build() models.UserImage {
	return models.UserImage{
		UserID:    u.UserID,
		Path:      u.Path,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
