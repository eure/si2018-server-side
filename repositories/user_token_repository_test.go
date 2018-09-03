package repositories

import (
	"testing"
	"github.com/tnkTaka/si2018-server-side/repositories"
	"fmt"
)

func TestUserTokenRepository_GetByToken(t *testing.T) {
	a := repositories.NewUserTokenRepository()
	token, err := a.GetByToken("USERTOKEN1")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(token)
}
