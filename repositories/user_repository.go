package repositories

import (
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/eure/si2018-server-side/entities"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(ent entities.User) error {
	s := engine.NewSession()
	if _, err := s.Insert(&ent); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(ent *entities.User) error {
	now := strfmt.DateTime(time.Now())

	s := engine.NewSession().Where("id = ?", ent.ID)
	ent.UpdatedAt = now
	if _, err := s.Update(ent); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByUserID(userID int64) (*entities.User, error) {
	var ent = entities.User{ID: userID}

	has, err := engine.Get(&ent)
	if err != nil {
		return nil, err
	}

	if has {
		return &ent, nil
	}

	return nil, nil
}

// limit / offset / 検索対象の性別 でユーザーを取得
// idsには取得対象に含めないUserIDを入れる (いいね/マッチ/ブロック済みなど)
func (r *UserRepository) FindWithCondition(limit, offset int, gender string, ids []int64) ([]entities.User, error) {
	var users []entities.User

	s := engine.NewSession()
	s.Where("gender = ?", gender)
	if len(ids) > 0 {
		s.NotIn("id", ids)
	}
	s.Limit(limit, offset)
	s.Desc("id")

	err := s.Find(&users)
	if err != nil {
		return users, err
	}

	return users, nil
}

// limit / offset / 検索対象の性別 でユーザーを取得
// idsには取得対象に含めないUserIDを入れる (いいね/マッチ/ブロック済みなど)
func (r *UserRepository) FindUsers(limit, offset int, gender string, ids []int64) ([]entities.User, error) {
	var users []entities.User

	s := engine.NewSession()

	s.Where("gender = ?", gender)
	if len(ids) > 0 {
		s.NotIn("id", ids)
	}
	s.Limit(limit, offset)
	s.Desc("created_at")

	err := s.Find(&users)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (r *UserRepository) FindByIDs(ids []int64) ([]entities.User, error) {
	var users []entities.User

	err := engine.In("id", ids).Find(&users)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (r *UserRepository) ParamsToUserEnt(u *entities.User, params si.PutProfileBody) *entities.User {
	user := &entities.User{
		ID:             u.ID,
		Nickname:       params.Nickname,
		ImageURI:       params.ImageURI,
		Tweet:          params.Tweet,
		Introduction:   params.Introduction,
		ResidenceState: params.ResidenceState,
		HomeState:      params.HomeState,
		Education:      params.Education,
		Job:            params.Job,
		AnnualIncome:   params.AnnualIncome,
		Height:         params.Height,
		BodyBuild:      params.BodyBuild,
		MaritalStatus:  params.MaritalStatus,
		Child:          params.Child,
		WhenMarry:      params.WhenMarry,
		WantChild:      params.WantChild,
		Smoking:        params.Smoking,
		Drinking:       params.Drinking,
		Holiday:        params.Holiday,
		HowToMeet:      params.HowToMeet,
		CostOfDate:     params.CostOfDate,
		NthChild:       params.NthChild,
		Housework:      params.Housework,
	}

	return user
}
