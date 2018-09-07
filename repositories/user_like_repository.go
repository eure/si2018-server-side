package repositories

import (
	"fmt"
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/strfmt"
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

func (r *UserLikeRepository) FindAllLikeUserResponse(ids []int64) ([]entities.LikeUserResponse, error) {
	// kadai4
	var usr []entities.User
	var likeresp []entities.LikeUserResponse

	s := engine.NewSession()
	err := s.In("id", ids).Find(&usr)
	if err != nil {
		return likeresp, err
	}
	likeresp[0].User.CreatedAt = usr[0].CreatedAt
	likeresp[0].LikedAt = strfmt.DateTime(time.Now())
	//fmt.Println(likeresp[0].User)
	//for i := 0; 0 < len(usr); i++ {
	//	likeresp[i].User = usr[i]
	//}
	fmt.Println(usr)
	return likeresp, nil
}
