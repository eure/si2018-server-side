package repositories

import (
	"github.com/go-xorm/builder"

	"github.com/eure/si2018-server-side/entities"
)

type UserMatchRepository struct{}

func NewUserMatchRepository() UserMatchRepository {
	return UserMatchRepository{}
}

func (r *UserMatchRepository) Create(ent entities.UserMatch) error {
	s := engine.NewSession()
	if _, err := s.Insert(&ent); err != nil {
		return err
	}

	return nil
}

func (r *UserMatchRepository) Get(userID, partnerID int64) (*entities.UserMatch, error) {
	var ent = entities.UserMatch{}
	var ids = []int64{userID, partnerID}
	has, err := engine.Where(builder.In("user_id", ids).And(builder.In("partner_id", ids))).Get(&ent)
	if err != nil {
		return nil, err
	}
	if has {
		return &ent, nil
	}
	return nil, nil
}

// マッチング済みのお相手一覧をlimit/offsetで取得する.
func (r *UserMatchRepository) FindByUserIDWithLimitOffset(userID int64, limit, offset int) ([]entities.UserMatch, error) {
	var matches []entities.UserMatch

	err := engine.Where("partner_id = ?", userID).Or("user_id = ?", userID).Limit(limit, offset).Desc("created_at").Find(&matches)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

// 自分が既にマッチングしている全てのお相手のUserIDを返す.
func (r *UserMatchRepository) FindAllByUserID(userID int64) ([]int64, error) {
	var matches []entities.UserMatch
	var ids []int64

	err := engine.Where("partner_id = ?", userID).Or("user_id = ?", userID).Find(&matches)
	if err != nil {
		return ids, err
	}

	for _, l := range matches {
		if l.UserID == userID {
			ids = append(ids, l.PartnerID)
			continue
		}
		ids = append(ids, l.UserID)
	}

	return ids, nil
}
