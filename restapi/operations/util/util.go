package util

import (
	"strings"
	"strconv"
	"errors"
	"fmt"
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

func ValidateToken(token string) error {
	// 文字列として正しいか
	if !isTokenStrValid(token) {
		return errors.New("Invalid token string ('USERTOKEN{id}' is required)")
	}
	// userが存在するか
	return nil
}

func isTokenStrValid(token string) bool {
	if !(strings.HasPrefix(token, "USERTOKEN")) {
		return false
	}

	id, err := strconv.Atoi(token[9:])
	if !(id > 0) || (err != nil) {
		return false
	}

	return true
}
/*
func userExists(token string) bool {
}
*/