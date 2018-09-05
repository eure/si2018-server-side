package util

import (
	"strings"
	"strconv"
	"errors"
	"fmt"
	"github.com/eure/si2018-server-side/repositories"
)

func GetIDByToken(token string) (int64, error) {
	r := repositories.NewUserTokenRepository()

	ut, err := r.GetByToken(token)
	if err != nil {
		return 0, err
	}

	if ut == nil {
		return 0, errors.New("Invalid token (user not exists)") // 必ずValidation後に呼ばれるならここいらない
	}

	return ut.UserID, nil
}

func ValidateToken(token string) error {
	// 文字列として正しいか
	if !isTokenStrValid(token) {
		return errors.New("Invalid token string ('USERTOKEN{id}' is required)")
	}
	// userが存在するか
	if !userExists(token) {
		return errors.New("Invalid token (user not exists)")
	}
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

func userExists(token string) bool {
	r := repositories.NewUserTokenRepository()
	ut, err := r.GetByToken(token)
	if err != nil {
		fmt.Println("Get id by token err:")
		fmt.Println(err)
		return false
	}

	if ut == nil {
		return false
	}

	return true
}