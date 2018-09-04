package user

import (
	"fmt"

	"github.com/eure/si2018-server-side/entities"
	"github.com/eure/si2018-server-side/models"
	"github.com/eure/si2018-server-side/repositories"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	nutr := repositories.NewUserTokenRepository()
	nulr := repositories.NewUserLikeRepository()
	nur := repositories.NewUserRepository()
	//ngur := repositories.NewGetUserRepository()
	// find userid
	usrid, err := nutr.GetByToken(p.Token)
	if err != nil {
		fmt.Println(err)
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "404",
				Message: "User Token Not Found",
			})
	}
	// find userlike
	userlike, err := nulr.FindLikeAll(usrid.UserID)
	if err != nil {
		fmt.Println(err)
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "404",
				Message: "UnknownFindfavorite",
			})
	}
	// find user
	userdesc, err := nur.GetByUserID(usrid.UserID)
	if err != nil {
		fmt.Println(err)
		return si.NewGetUsersInternalServerError().WithPayload(
			&si.GetUsersInternalServerErrorBody{
				Code:    "404",
				Message: "User Token Not Found",
			})
	}
	ent, err := nur.FindWithCondition(int(p.Limit), int(p.Offset), userdesc.GetOppositeGender(), userlike)
	if err != nil {
		fmt.Println(err)
		si.NewGetUsersInternalServerError()
	}

	var ud []*models.User
	for i := 0; i < len(ent); i++ {
		us := ent[i].Build()
		ud = append(ud, &us)
	}

	return si.NewGetUsersOK().WithPayload(ud)

}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	nur := repositories.NewUserRepository()
	nut := repositories.NewUserTokenRepository()
	usrid, err := nur.GetByUserID(p.UserID)
	if err != nil {
		fmt.Println(err)
	}
	usrtoken, err := nut.GetByToken(p.Token)
	if err != nil {
		fmt.Println(err)
	}
	if usrid.ID != usrtoken.UserID {
		fmt.Println("error")
	}
	sEnt := usrid.Build()
	return si.NewGetProfileByUserIDOK().WithPayload(&sEnt)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	nur := repositories.NewUserRepository()
	user, err := nur.GetByUserID(p.UserID)
	if err != nil {
		fmt.Println(err)
	}
	// token validation implemtation later

	/*
		updateuser := reflect.Indirect(reflect.ValueOf(p.Params))
		updatefield := updateuser.Type()
		usr := reflect.Indirect(reflect.ValueOf(user))
		userfield := usr.Type()
		for i := 0; i < updatefield.NumField(); i++ {
			nic := updatefield.Field(i).Name
			up := updateuser.FieldByName(nic).Interface()
			//if up == nil {
			//	break
			//}
			fmt.Println(up.(string))
			if updatefield.Field(i).Name == userfield.Field(i).Name {

				user.Nickname = up.(string)
				fmt.Println(up.(string))
				break
			}
		}
	*/
	//for i := 0; i < userfield.NumField(); i++ {
	//	fmt.Println(userfield.Field(i).Name)
	//
	//}
	binduser(p.Params, user)
	err = nur.Update(user)
	if err != nil {
		fmt.Println(err)
	}
	respuser, err := nur.GetByUserID(p.UserID)
	if err != nil {
		fmt.Println(err)
	}
	updateuser := respuser.Build()
	return si.NewPutProfileOK().WithPayload(&updateuser)
}

func binduser(user si.PutProfileBody, ent *entities.User) {
	ent.AnnualIncome = user.AnnualIncome
	ent.BodyBuild = user.BodyBuild
	ent.Child = user.Child
	ent.CostOfDate = user.CostOfDate
	ent.Drinking = user.Drinking
	ent.Education = user.Education
	ent.Height = user.Height
	ent.Holiday = user.Holiday
	ent.HomeState = user.HomeState
	ent.Housework = user.Housework
	ent.HowToMeet = user.HowToMeet
	ent.ImageURI = user.ImageURI
	ent.Introduction = user.Introduction
	ent.Job = user.Job
	ent.MaritalStatus = user.MaritalStatus
	ent.Nickname = user.Nickname
	ent.NthChild = user.NthChild
	ent.ResidenceState = user.ResidenceState
	ent.Smoking = user.Smoking
	ent.Tweet = user.Tweet
	ent.WantChild = user.WantChild
	ent.WhenMarry = user.WhenMarry

}
