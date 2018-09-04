package user

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/restapi/operations/util"
	"fmt"
	"reflect"
)

func GetUsers(p si.GetUsersParams) middleware.Responder {
	limit := p.Limit
	offset := p.Offset
	token := p.Token
	/* TODO range validation */
	/* TODO token validation */
	
	ru := repositories.NewUserRepository()
	//rt := repositories.NewUserTokenRepository()
	rl := repositories.NewUserLikeRepository()

	/*ut, err := rt.GetByToken(token)
	id := ut.UserID
	*/
	id, err := util.GetIDByToken(token)
	user, err := ru.GetByUserID(id)
	likes, err := rl.FindLikeAll(id)

	users_ent, err := ru.FindWithCondition(int(limit), int(offset), user.GetOppositeGender(), likes);
	if err != nil {
		fmt.Println(err)
		return si.NewGetUsersInternalServerError()
	}

	var users entities.Users
	users = entities.Users(users_ent)
	return si.NewGetUsersOK().WithPayload(users.Build())
}

func GetProfileByUserID(p si.GetProfileByUserIDParams) middleware.Responder {
	id := p.UserID
	//token := p.Token
	/* TODO token validation */

	if id < 1 { /* TODO */
		return si.NewGetProfileByUserIDBadRequest().WithPayload(
			&si.GetProfileByUserIDBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	r := repositories.NewUserRepository()

	user, err := r.GetByUserID(id)
	if err != nil {
		fmt.Println(err)
		return si.NewGetUsersInternalServerError()
	}
	
	res := user.Build()

	return si.NewGetProfileByUserIDOK().WithPayload(&res)
}

func PutProfile(p si.PutProfileParams) middleware.Responder {
	id := p.UserID
	ps := p.Params
	//token := ps.Token
	/* TODO token validation */
	r := repositories.NewUserRepository()
	user, _ := r.GetByUserID(id)

	//fmt.Printf("user:%+v\n", user)
	//fmt.Printf("param:%+v\n", p.Params)

	user.AnnualIncome = ps.AnnualIncome
	user.BodyBuild = ps.BodyBuild
	user.Child = ps.Child
	user.CostOfDate = ps.CostOfDate
	user.Drinking = ps.Drinking
	user.Education = ps.Education
	user.Height = ps.Height
	user.Holiday = ps.Holiday
	user.HomeState = ps.HomeState
	user.Housework = ps.Housework
	user.HowToMeet = ps.HowToMeet
	user.ImageURI = ps.ImageURI
	user.Introduction = ps.Introduction
	user.Job = ps.Job
	user.MaritalStatus = ps.MaritalStatus
	user.Nickname = ps.Nickname
	user.NthChild = ps.NthChild
	user.ResidenceState = ps.ResidenceState
	user.Smoking = ps.Smoking
	user.Tweet = ps.Tweet
	user.WantChild = ps.WantChild
	user.WhenMarry = ps.WhenMarry

	//fmt.Printf("newu:%+v\n", user)
	err := r.Update(user)
	fmt.Println(err) /* TODO err */
	//fmt.Println("UPDATE DONE")
	//fmt.Printf("user:%+v\n", user)

	/*v := reflect.ValueOf(params)
	fmt.Println(v)

    values := make([]interface{}, v.NumField())

    for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		
	}*/
	
	/*if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.BodyBuild != "" {
		user.BodyBuild = params.BodyBuild
	}
	if params.Child != "" {
		user.Child = params.Child
	}
	if params.CostOfDate != "" {
		user.CostOfDate = params.CostOfDate
	}
	if params.Drinking != "" {
		user.Drinking = params.Drinking
	}
	if params.Education != "" {
		user.Education = params.Education
	}
	if params.Height != "" {
		user.Height = params.Height
	}
	if params.Holiday != "" {
		user.Holiday = params.Holiday
	}
	if params.HomeState != "" {
		user.HomeState = params.HomeState
	}
	if params.Housework != "" {
		user.Housework = params.Housework
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}
	if params.AnnualIncome != "" {
		user.AnnualIncome = params.AnnualIncome
	}*/



	return si.NewPutProfileOK()
}

func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()
  
	for i := 0; i < size; i++ {
	  field := elem.Type().Field(i).Name
	  value := elem.Field(i).Interface()
	  result[field] = value
	}
  
	return result
  }