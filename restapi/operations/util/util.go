package util

import (
	"github.com/eure/si2018-server-side/repositories"
)

func GetIdByToken(token string) int64, error {
	rt := repositories.NewUserTokenRepository()

	ut, err := rt.GetByToken(token)
	if err != nil {
		return nil, error
	}

	return ut.UserID, nil
}


