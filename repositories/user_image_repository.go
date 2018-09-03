package repositories

import (
	"time"

	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/strfmt"
)

type UserImageRepository struct{}

func NewUserImageRepository() UserImageRepository {
	return UserImageRepository{}
}

func (r *UserImageRepository) Create(ent entities.UserImage) error {
	s := engine.NewSession()
	if _, err := s.Insert(&ent); err != nil {
		return err
	}

	return nil
}

func (r *UserImageRepository) Update(ent entities.UserImage) error {
	now := strfmt.DateTime(time.Now())
	s := engine.NewSession().Where("user_id = ?", ent.UserID)
	ent.UpdatedAt = now

	if _, err := s.Update(ent); err != nil {
		return err
	}
	return nil
}

func (r *UserImageRepository) GetByUserID(userID int64) (*entities.UserImage, error) {
	var ent = entities.UserImage{UserID: userID}

	has, err := engine.Get(&ent)
	if err != nil {
		return nil, err
	}

	if has {
		return &ent, nil
	}

	return nil, nil
}

func (r *UserImageRepository) GetByUserIDs(userIDs []int64) ([]entities.UserImage, error) {
	var userImages []entities.UserImage

	err := engine.In("user_id", userIDs).Find(&userImages)
	if err != nil {
		return userImages, err
	}

	return userImages, nil
}
