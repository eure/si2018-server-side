package repositories

import (
	"fmt"
	"log"

	"github.com/eure/si2018-server-side/entities"
	"github.com/go-xorm/builder"
)

type UserLikeRepository struct{}

func NewUserLikeRepository() UserLikeRepository {
	return UserLikeRepository{}
}

func (r *UserLikeRepository) Create(ent entities.UserLike) error {
	s := engine.NewSession()
	if _, err := s.Insert(&ent); err != nil {
		return err
	}

	return nil
}

func (r *UserLikeRepository) Get(userID, partnerID int64) (*entities.UserLike, error) {
	var like entities.UserLike
	var ids = []int64{userID, partnerID}
	s := engine.NewSession()
	defer func() { log.Println(s.LastSQL()) }()
	has, err := s.Where(builder.In("user_id", ids).And(builder.In("partner_id", ids))).Get(&like)
	if err != nil {
		return nil, err
	}
	if has {
		return &like, nil
	}
	return nil, nil
}

func (r *UserLikeRepository) GetLikeByUserID(meID, partnerID int64) (*entities.UserLike, error) {
	var like entities.UserLike

	_, err := engine.Where("user_id = ?", partnerID).And("partner_id = ?", meID).Get(&like)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &like, nil
}

func (r *UserLikeRepository) FindGotLikeWithLimitOffset(userID int64, limit, offset int, matchIDs []int64) ([]entities.UserLike, error) {
	var likes []entities.UserLike

	s := engine.NewSession()
	s.Where("partner_id = ?", userID)
	if len(matchIDs) > 0 {
		s.NotIn("user_id", matchIDs)
	}
	s.Limit(limit, offset)
	s.Desc("created_at")
	err := s.Find(&likes)
	if err != nil {
		return likes, err
	}

	return likes, nil
}

func (r *UserLikeRepository) FindSendLikeWithLimitOffset(userID int64, limit, offset int) ([]entities.UserLike, error) {
	var likes []entities.UserLike

	err := engine.Where("user_id = ?", userID).Limit(limit, offset).Desc("created_at").Find(&likes)
	if err != nil {
		return likes, err
	}

	return likes, nil
}

// Likeした・された全てのUserIDを返す
func (r *UserLikeRepository) FindLikeAll(userID int64) ([]int64, error) {
	var likes []entities.UserLike
	var ids []int64

	err := engine.Where("partner_id = ?", userID).Or("user_id = ?", userID).Find(&likes)
	if err != nil {
		return ids, err
	}

	for _, l := range likes {
		if l.UserID == userID {
			ids = append(ids, l.PartnerID)
			continue
		}
		ids = append(ids, l.UserID)
	}

	return ids, nil
}
