package entities

import (
	"github.com/eure/si2018-server-side/models"
	"github.com/go-openapi/strfmt"
	"time"
)

type UserImage struct {
	UserID    int64           `xorm:"user_id"`
	Path      string          `xorm:"path"`
	CreatedAt strfmt.DateTime `xorm:"created_at"`
	UpdatedAt strfmt.DateTime `xorm:"updated_at"`
}

func NewUserImage(userID int64, path string) UserImage {
	now := strfmt.DateTime(time.Now())
	return UserImage{
		UserID: userID,
		Path: path,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u UserImage) Build() models.UserImage {
	return models.UserImage{
		UserID:    u.UserID,
		Path:      u.Path,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
