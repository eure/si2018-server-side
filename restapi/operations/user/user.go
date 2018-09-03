package user

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	return si.NewGetUsersOK()
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	return si.NewGetProfileByUserIDOK()
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	return si.NewPutProfileOK()
}

func getUsersInteralServerErrorResponse(message string) middleware.Responder {
  return si.NewGetUsersInternalServerError().WithPayload(
    &si.GetUsersInternalServerErrorBody {
      Code:    "500",
      Message: message,
  })
}

func getUsersUnauthorizedResponse(message string) middleware.Responder {
  return si.NewGetUsersUnauthorized().WithPayload(
    &si.GetUsersUnauthorizedBody{
      Code:    "401",
      Message: "Your Token Is Invalid",
    })
  }
