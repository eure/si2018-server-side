package repositories

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"os"

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
	s := engine.NewSession().Where("user_id = ?", ent.UserID)

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

func (r *UserImageRepository) SaveImageAssets(base64Str strfmt.Base64, userID int64, typeName string) (string, error) {
	imagePath := fmt.Sprintf("%d.%s", userID, typeName)
	file, err := os.Create("./assets/" + imagePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = file.Write(base64Str)
	if err != nil {
		return "", err
	}

	return imagePath, err
}

// formatを返しているのにvalidationはおかしい気がする
func (r *UserImageRepository) ImageValidation(base64Str strfmt.Base64) (string, error) {
	reader := bytes.NewReader(base64Str)
	imageByteSize := reader.Size()

	if imageByteSize > 5e6 {
		return "", errors.New("画像サイズが大きいです(5MBまで)")
	}

	_, typeName, err := image.DecodeConfig(reader)
	if err != nil {
		return "", err
	}

	if typeName != "jpeg" && typeName != "png" {
		return "", errors.New("画像フォーマットが適切ではありません")
	}

	return typeName, nil
}
