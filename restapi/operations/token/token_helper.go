package token

import (
	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/repositories"
)

func GetUserByToken(token string) (*entities.User, error) {
	ent, err := repositories.NewUserTokenRepository().GetByToken(token)
	if err != nil {
		return nil, err
	}
	if ent == nil {
		return nil, nil
	}

	user, err := repositories.NewUserRepository().GetByUserID(ent.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	return user, nil
}
