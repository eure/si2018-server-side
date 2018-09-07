package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/restapi/operations/util"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"fmt"
	"time"
	"github.com/go-openapi/strfmt"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	limit := p.Limit
	//limit := 20
	offset := p.Offset
	token := p.Token
	
	// Validations
	err1 := util.ValidateLimit(limit)
	err2 := util.ValidateOffset(offset)
	if (err1 != nil) || (err2 != nil) {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	err := util.ValidateToken(token)
	if err != nil {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	ru := repositories.NewUserRepository()
	rm := repositories.NewUserMatchRepository()
	rl := repositories.NewUserLikeRepository()

	id, _ := util.GetIDByToken(token) // ValidateToken()で先に呼ばれているのでerrは潰してよい

	// Get users already matching
	matches, err := rm.FindAllByUserID(id)
	if err != nil {
		fmt.Print("Find matches err: ") /* TODO use log */
		fmt.Println(err)
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// Get a list of UserLike the user got 
	likes, err := rl.FindGotLikeWithLimitOffset(id, int(limit), int(offset), matches)
	if err != nil {
		fmt.Print("Find likes err: ")
		fmt.Println(err)
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}

	// レスポンス作成ループ内でDB問い合わせを繰り返さないためにidsとumを準備
	ids := make([]int64, 0) /* TODO can use map's key as ids slice? */
	for _, l := range likes {
		ids = append(ids, l.UserID)
	}

	users, err := ru.FindByIDs(ids)
	if err != nil {
		fmt.Print("Find users by ids err: ")
		fmt.Println(err)
		return si.NewGetLikesInternalServerError().WithPayload(
			&si.GetLikesInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	um := make(map[int64]entities.User)
	for _, u := range users {
		//u := u
		um[u.ID] = u
	}

	res := make([]entities.LikeUserResponse, 0)
	
	for _, l := range likes {
		r := entities.LikeUserResponse{}
		r.LikedAt = l.CreatedAt
		r.ApplyUser(um[l.UserID])
		res = append(res, r)
	}

	var reses entities.LikeUserResponses
	reses = res

	return si.NewGetLikesOK().WithPayload(reses.Build())
}

func PostLike(p si.PostLikeParams) middleware.Responder {
	token := p.Params.Token
	rid := p.UserID

	err := util.ValidateToken(token)
	if err != nil {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "401",
				Message: "Token Is Invalid",
			})
	}

	sid, _ := util.GetIDByToken(token)

	ru := repositories.NewUserRepository()
	rl := repositories.NewUserLikeRepository()

	// Same gender?
	users, err := ru.FindByIDs([]int64{rid, sid})
	if err != nil {
		fmt.Print("FindByIDs err: ")
		fmt.Println(err)
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if len(users) != 2 || users[0].Gender == users[1].Gender { // Believe short-circuit evaluation
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request (Same gender)",
			})
	}

	ul, err := rl.GetLikeBySenderIDReceiverID(sid, rid)
	if err != nil {
		fmt.Print("Get like by ids (first) err: ")
		fmt.Println(err)
		return si.NewPostLikeInternalServerError().WithPayload(
			&si.PostLikeInternalServerErrorBody{
				Code:    "500",
				Message: "Internal Server Error",
			})
	}
	if ul == nil { 
		ul, err := rl.GetLikeBySenderIDReceiverID(rid, sid)
		if err != nil {
			fmt.Print("Get like by ids (second) err: ")
			fmt.Println(err)
			return si.NewPostLikeInternalServerError().WithPayload(
				&si.PostLikeInternalServerErrorBody{
					Code:    "500",
					Message: "Internal Server Error",
				})
		}

		if ul == nil { // If the first like
			l := entities.UserLike{}
		    l.UserID = sid
		    l.PartnerID = rid
		    l.CreatedAt = strfmt.DateTime(time.Now())
		    l.UpdatedAt = strfmt.DateTime(time.Now())
		    err := rl.Create(l)
		    if err != nil {
		    	fmt.Print("Create first like err: ")
				fmt.Println(err)
				return si.NewPostLikeInternalServerError().WithPayload(
					&si.PostLikeInternalServerErrorBody{
						Code:    "500",
						Message: "Internal Server Error",
					})
		    }
		} else { // If match
			l := entities.UserLike{}
		    l.UserID = sid
		    l.PartnerID = rid
		    l.CreatedAt = strfmt.DateTime(time.Now())
		    l.UpdatedAt = strfmt.DateTime(time.Now())
		    err := rl.Create(l)
		    if err != nil {
		    	fmt.Print("Create matching like err: ")
				fmt.Println(err)
				return si.NewPostLikeInternalServerError().WithPayload(
					&si.PostLikeInternalServerErrorBody{
						Code:    "500",
						Message: "Internal Server Error",
					})
		    }
    
		    m := entities.UserMatch{}
		    m.UserID = rid
		    m.PartnerID = sid
		    m.CreatedAt = strfmt.DateTime(time.Now())
		    m.UpdatedAt = strfmt.DateTime(time.Now())
    
		    rm := repositories.NewUserMatchRepository()
		    err = rm.Create(m)
		    if err != nil {
				fmt.Print("Create match err: ")
				fmt.Println(err)
				return si.NewPostLikeInternalServerError().WithPayload(
					&si.PostLikeInternalServerErrorBody{
						Code:    "500",
						Message: "Internal Server Error",
					})
			}
		}
	} else { // If duplicate
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request (Duplicate like)",
			})
	}

	return si.NewPostLikeOK()
}