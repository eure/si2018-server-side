package repositories

import (
	"log"

	"github.com/go-openapi/strfmt"
	"github.com/go-xorm/builder"

	"time"

	"github.com/eure/si2018-server-side/entities"
)

type UserMessageRepository struct{}

func NewUserMessageRepository() UserMessageRepository {
	return UserMessageRepository{}
}

func (r *UserMessageRepository) Create(ent entities.UserMessage) error {
	now := strfmt.DateTime(time.Now())

	s := engine.NewSession()
	ent.CreatedAt = now
	ent.UpdatedAt = now
	if _, err := s.Insert(&ent); err != nil {
		return err
	}

	return nil
}

//　最後に送信したメッセージを取得する。
func (r *UserMessageRepository) GetLastMessages(id int64) (*entities.UserMessage, error) {
	var ent entities.UserMessage

	has, err := engine.Where("user_id = ?", id).Desc("created_at").Get(&ent)
	if err != nil {
		return nil, err
	}

	if has {
		return &ent, nil
	}

	return nil, nil
}

// userとpartnerがやりとりしたメッセージをlimit/latest/oldestで取得する.
func (r *UserMessageRepository) GetMessages(userID, partnerID int64, limit int, latest, oldest *strfmt.DateTime) ([]entities.UserMessage, error) {
	var messages []entities.UserMessage
	var ids = []int64{userID, partnerID}

	s := engine.NewSession()
	defer func() { log.Println(s.LastSQL()) }()
	s.Where(builder.In("user_id", ids))
	s.And(builder.In("partner_id", ids))
	if latest != nil {
		s.And("created_at < ?", latest)
	}
	if oldest != nil {
		s.And("created_at > ?", oldest)
	}
	s.Limit(limit)
	s.Desc("created_at")
	err := s.Find(&messages)
	if err != nil {
		return messages, err
	}

	return messages, nil
}
