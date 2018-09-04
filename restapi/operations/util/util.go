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
/*
func ValidateToken(token string) error {
	// 文字列として正しいか
	// userが存在するか
}

func isTokenStrValid(token string) bool {
}


*/