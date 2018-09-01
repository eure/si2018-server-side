package repositories

import (
	"github.com/eure/si2018-server-side/entities"
)

type UserTokenRepository struct{}

func NewUserTokenRepository() *UserTokenRepository {
	return &UserTokenRepository{}
}

func (r *UserTokenRepository) Create(ent entities.UserToken) error {
	s := engine.NewSession()
	if _, err := s.Insert(&ent); err != nil {
		return err
	}

	return nil
}

func (r *UserTokenRepository) Update(ent entities.UserToken, cols []string) error {
	s := engine.NewSession()
	s.MustCols(cols...)
	if _, err := s.Update(ent); err != nil {
		return err
	}
	return nil
}

func (r *UserTokenRepository) GetByUserID(userID int64) (*entities.UserToken, error) {
	var ent = entities.UserToken{UserID: userID}

	has, err := engine.Get(&ent)
	if err != nil {
		return nil, err
	}

	if has {
		return &ent, nil
	}

	return nil, nil
}

func (r *UserTokenRepository) GetByToken(token string) (*entities.UserToken, error) {
	var ent = entities.UserToken{Token: token}

	has, err := engine.Get(&ent)
	if err != nil {
		return nil, err
	}

	if has {
		return &ent, nil
	}

	return nil, nil
}
