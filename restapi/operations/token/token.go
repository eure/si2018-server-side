package token

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

// DB アクセス: 1 回
// 計算量: O(1)
func GetTokenByUserID(p si.GetTokenByUserIDParams) middleware.Responder {
	tokenRepo := repositories.NewUserTokenRepository()

	token, err := tokenRepo.GetByUserID(p.UserID)
	if err != nil {
		return si.GetTokenByUserIDThrowInternalServerError(err)
	}
	if token == nil {
		return si.GetTokenByUserIDThrowNotFound()
	}

	sEnt := token.Build()
	return si.NewGetTokenByUserIDOK().WithPayload(&sEnt)
}
