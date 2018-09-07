package entities

import (
	"github.com/eure/si2018-server-side/models"
	"github.com/go-openapi/strfmt"
	"time"
)

type UserMessage struct {
	UserID    int64           `xorm:"user_id"`
	PartnerID int64           `xorm:"partner_id"`
	Message   string          `xorm:"message"`
	CreatedAt strfmt.DateTime `xorm:"created_at"`
	UpdatedAt strfmt.DateTime `xorm:"updated_at"`
}

func NewUserMessage(userID int64, partnerID int64, msg string) UserMessage {
	return UserMessage{
		UserID: userID,
		PartnerID: partnerID,
		Message: msg,
		CreatedAt: strfmt.DateTime(time.Now()),
		UpdatedAt: strfmt.DateTime(time.Now()),
	}
}

func (u UserMessage) Build() models.UserMessage {
	return models.UserMessage{
		UserID:    u.UserID,
		PartnerID: u.PartnerID,
		Message:   u.Message,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

type UserMessages []UserMessage

func (msgs *UserMessages) Build() []*models.UserMessage {
	var sMsgs []*models.UserMessage

	for _, m := range *msgs {
		sMsg := m.Build()
		sMsgs = append(sMsgs, &sMsg)
	}
	return sMsgs
}
