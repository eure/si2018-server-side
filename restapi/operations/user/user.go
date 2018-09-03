package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/entities"
  "github.com/eure/si2018-server-side/repositories"
  si "github.com/eure/si2018-server-side/restapi/summerintern"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
  //FIXME Token認証は共通化したい
  r := repositories.NewUserTokenRepository()

  ent, err := r.GetByToken(p.Token)

  if err != nil { return getUsersInteralServerErrorResponse("Internal Server Error") }
  if ent == nil { return getUsersUnauthorizedResponse("Your Token Is Invalid") }

  userToken := ent.Build()
  userID := userToken.UserID

  userLikeRepository := repositories.NewUserLikeRepository()
  exclusionIds, err := userLikeRepository.FindLikeAll(userID)

  if err != nil { return getUsersInteralServerErrorResponse("Internal Server Error") }

  userRepository := repositories.NewUserRepository()
  userEnt, err  := userRepository.GetByUserID(userID)
  if err != nil { return getUsersInteralServerErrorResponse("Internal Server Error") }

  // int64になっているのでcastする必要がある
  limit := int(p.Limit)
  offset := int(p.Offset)
  gender := userEnt.GetOppositeGender()

  var usersEnt entities.Users
  usersEnt, err = userRepository.FindUsers(limit, offset, gender, exclusionIds)
  if err != nil { return getUsersInteralServerErrorResponse("Internal Server Error") }

  users := usersEnt.Build()
  return si.NewGetUsersOK().WithPayload(users)
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
