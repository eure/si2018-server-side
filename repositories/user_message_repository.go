package repositories

import (
	"errors"
	"log"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-xorm/builder"

	"github.com/eure/si2018-server-side/entities"
)

type UserMessageRepository struct{}

func NewUserMessageRepository() UserMessageRepository {
	return UserMessageRepository{}
}

func (r *UserMessageRepository) Create(ent entities.UserMessage) error {
	now := strfmt.DateTime(time.Now())
	ent.CreatedAt = now
	ent.UpdatedAt = now

	s := engine.NewSession()
	if _, err := s.Insert(&ent); err != nil {
		return err
	}

	return nil
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
		s.And("created_at > ?", latest)
	}
	if oldest != nil {
		s.And("created_at < ?", oldest)
	}
	s.Limit(limit)
	err := s.Find(&messages)
	if err != nil {
		return messages, err
	}

	return messages, nil
}

func (r *UserMessageRepository) Validate(u entities.UserMessage) []error {
	var res []error
	if err := isMessagePresence(u.Message); err != nil {
		res = append(res, err)
	}

	if err := isMessageLength(u.Message); err != nil {
		res = append(res, err)
	}

	if err := isMatched(u); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return res
	}

	return nil
}

func isMessagePresence(message string) error {
	if len(message) == 0 {
		return errors.New("メッセージ内容を入力してください")
	}

	return nil
}

func isMessageLength(message string) error {
	if len(message) >= 5000 {
		return errors.New("最大5000文字まで送信できます")
	}

	return nil
}

func isMatched(u entities.UserMessage) error {
	var matches []entities.UserMatch

	engine.
		Where("partner_id = ?", u.UserID).And("user_id = ?", u.PartnerID).
		Or("partner_id = ?", u.PartnerID).And("user_id = ?", u.UserID).
		Find(&matches)

	if len(matches) == 0 {
		return errors.New("マッチング済みの相手にしかメッセージを送信できません")
	}

	return nil
}
