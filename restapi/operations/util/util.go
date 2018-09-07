package util

import (
	"errors"
	"log"
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
	if token == ""{
		return errors.New("Invalid token (empty token)")
	}
	if !userExists(token) {
		return errors.New("Invalid token (user not exists)")
	}
	return nil
}

func userExists(token string) bool {
	r := repositories.NewUserTokenRepository()
	ut, err := r.GetByToken(token)
	if err != nil {
		log.Print("Get id by token err:", err)
		return false
	}

	if ut == nil {
		return false
	}

	return true
}

func ValidateLimit(lim int64) error {
	if lim <= 0 { // If Limit 0, DB returns 0 result (maybe) <= NO!!!!!!!
		return errors.New("Limit must be >= 0")
	}
	return nil
}

func ValidateOffset(ofs int64) error {
	if ofs < 0 {
		return errors.New("Offset must be >= 0")
	}
	return nil
}
