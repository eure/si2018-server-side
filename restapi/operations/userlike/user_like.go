package userlike

import (
	"github.com/go-openapi/runtime/middleware"

	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"fmt"
)

func GetLikes(p si.GetLikesParams) middleware.Responder {
	//limit := int(p.Limit)
	limit := 20
	offset := int(p.Offset)
	//token := p.Token
	/* TODO token validation */

	ru := repositories.NewUserRepository()
	rm := repositories.NewUserMatchRepository()
	rl := repositories.NewUserLikeRepository()

	/* TODO get id by token */
	//id := GetIdByToken() //int64
	id := int64(1111)

	matches, err := rm.FindAllByUserID(id)
	if err != nil {
		fmt.Print("Find matches err: ")
		fmt.Println(err)
	}

	likes, err := rl.FindGotLikeWithLimitOffset(id, limit, offset, matches)
	if err != nil {
		fmt.Print("Find likes err: ")
		fmt.Println(err)
	}

	ids := make([]int64, 0) /* TODO can use map's key as ids slice? */
	for _, l := range likes {
		ids = append(ids, l.UserID)
	}

	users, err := ru.FindByIDs(ids)
	if err != nil {
		fmt.Print("Find users by ids err: ")
		fmt.Println(err)
	}
	um := make(map[int64]entities.User)
	for _, u := range users {
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
	/*token := p.Params.Token
	rid := p.UserID
	sid := GetIdByToken(token)

	r := NewUserLikeRepository()
	like, err := r.GetLikeBySenderIDReceiverID(sid, rid)
	if like == nil { // If the first like
        /* TODO UPDATE? */
	//}
	/* TODO same gender */
	/* TODO dup */

	return si.NewPostLikeOK()
}