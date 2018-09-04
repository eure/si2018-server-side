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
	//limit := int(p.Limit)
	limit := 20
	offset := int(p.Offset)
	token := p.Token
	/* TODO token validation */

	ru := repositories.NewUserRepository()
	rm := repositories.NewUserMatchRepository()
	rl := repositories.NewUserLikeRepository()

	id, err := util.GetIDByToken(token)
	//id := int64(2)
	if err != nil {
		fmt.Print("Get id err: ")
		fmt.Println(err)
	}

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
	fmt.Println(likes)

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
	//token := p.Params.Token
	rid := p.UserID
	//sid := GetIdByToken(token)
	sid := int64(1)

	ru := repositories.NewUserRepository()
	rl := repositories.NewUserLikeRepository()

	// Not same id?
	if rid == sid {
		/* TODO */
	}

	// Same gender?
	users, err := ru.FindByIDs([]int64{rid, sid})
	if err != nil {
		fmt.Print("FindByIDs err: ")
		fmt.Println(err)
	}
	if users[0].Gender == users[1].Gender {
		fmt.Println("Same gender err")
		
	}

	ul, err := rl.GetLikeBySenderIDReceiverID(sid, rid)
	if ul == nil { 
		ul, err := rl.GetLikeBySenderIDReceiverID(rid, sid)
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
		    
			}
		}
	} else {// If duplicate
		
	}
		
	/*ul, err := rl.GetLikeBySenderIDReceiverID(sid, rid)
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
		} 
	} else if ul.UserID == sid { // If duplicate
	} else if ul.PartnerID == sid { // If match
		l := entities.UserLike{}
		l.UserID = sid
		l.PartnerID = rid
		l.CreatedAt = strfmt.DateTime(time.Now())
		l.UpdatedAt = strfmt.DateTime(time.Now())
		err := rl.Create(l)
		if err != nil {
			fmt.Print("Create matching like err: ")
			fmt.Println(err)
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
		}
	} /* TODO check already matching? */

	return si.NewPostLikeOK()
}