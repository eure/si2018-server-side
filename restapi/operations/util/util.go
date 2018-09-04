package util

import (
	"github.com/eure/si2018-server-side/repositories"
)

func GetIDByToken(token string) (int64, error) {
	rt := repositories.NewUserTokenRepository()

	ut, err := rt.GetByToken(token)
	if err != nil {
		return int64(-1), err
	}

	return ut.UserID, nil
}


