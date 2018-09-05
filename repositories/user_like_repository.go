package repositories

import (
	"errors"
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/strfmt"
)

type UserLikeRepository struct{}

func NewUserLikeRepository() UserLikeRepository {
	return UserLikeRepository{}
}

func (r *UserLikeRepository) Create(ent entities.UserLike) error {
	now := strfmt.DateTime(time.Now())
	ent.CreatedAt = now
	ent.UpdatedAt = now

	s := engine.NewSession()

	if err := s.Begin(); err != nil {
		return err
	}

	if _, err := s.Insert(&ent); err != nil {
		s.Rollback()
		return err
	}

	var user entities.UserLike
	has, err := engine.
		Where("user_id = ?", ent.PartnerID).
		And("partner_id = ?", ent.UserID).
		Get(&user)

	if err != nil {
		s.Rollback()
		return err
	}

	if has {
		// マッチングさせる
		now = strfmt.DateTime(time.Now())
		userMatch := entities.UserMatch{
			UserID:    ent.UserID,
			PartnerID: ent.PartnerID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if _, err := s.Insert(&userMatch); err != nil {
			s.Rollback()
			return err
		}
	}

	return s.Commit()
}

// 自分が既にLikeしている/されている状態の全てのUserのIDを返す.
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

// いいねを1件取得する.
// userIDはいいねを送った人, partnerIDはいいねを受け取った人.
func (r *UserLikeRepository) GetLikeBySenderIDReceiverID(userID, partnerID int64) (*entities.UserLike, error) {
	var ent entities.UserLike

	has, err := engine.Where("user_id = ?", userID).And("partner_id = ?", partnerID).Get(&ent)
	if err != nil {
		return nil, err
	}
	if has {
		return &ent, nil
	}
	return nil, nil
}

// マッチ済みのお相手を除き、もらったいいねを、limit/offsetで取得する.
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

func (r *UserLikeRepository) Validate(userLike entities.UserLike) []error {
	var res []error

	if err := isHeterosexual(userLike); err != nil {
		res = append(res, err)
	}

	if err := isLiked(userLike); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return res
	}

	return nil
}

func isLiked(userLike entities.UserLike) error {
	var ent entities.UserLike

	has, _ := engine.
		Where("user_id = ?", userLike.UserID).
		And("partner_id = ?", userLike.PartnerID).
		Get(&ent)

	if has {
		return errors.New("すでにいいねしています")
	}

	return nil
}

func isHeterosexual(userLike entities.UserLike) error {
	var sender = entities.User{ID: userLike.UserID}

	has, _ := engine.Get(&sender)

	if !has {
		return errors.New("送信ユーザーが見つかりません")
	}

	var destination = entities.User{ID: userLike.PartnerID}

	has, _ = engine.Get(&destination)

	if !has {
		return errors.New("送信先ユーザーが見つかりません")
	}

	// MEMO 同性へのいいねはできてもいいような気がするがどうだろう？
	if sender.Gender == destination.Gender {
		return errors.New("同性へのいいねはできません")
	}

	return nil
}
