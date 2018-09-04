package message

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/eure/si2018-server-side/restapi/operations/util"
	"github.com/eure/si2018-server-side/repositories"
	"github.com/eure/si2018-server-side/entities"
	"github.com/go-openapi/runtime/middleware"
	"fmt"
	"time"
	"github.com/go-openapi/strfmt"
)

func PostMessage(p si.PostMessageParams) middleware.Responder {
	message := p.Params.Message
	token := p.Params.Token
	pid := p.UserID

	err := util.ValidateToken(token)
	if err != nil {
		fmt.Println("Invalid token err:")
		fmt.Println(err)
	}

	uid, err := util.GetIDByToken(token)
	if err != nil {
		fmt.Print("Get id err: ")
		fmt.Println(err)
	}

	rl := repositories.NewUserLikeRepository()
	like, err := rl.GetLikeBySenderIDReceiverID(uid, pid)
	if err != nil {
		fmt.Print("Get likes err: ")
		fmt.Println(err)
	}
	if like == nil {
		fmt.Print("Not matching yet")
		
	}

	ent := entities.UserMessage{}
	ent.UserID = uid
	ent.PartnerID = pid
	ent.Message = message
	ent.CreatedAt = strfmt.DateTime(time.Now())
	ent.UpdatedAt = strfmt.DateTime(time.Now())

	rm := repositories.NewUserMessageRepository()
	err = rm.Create(ent)
	if err != nil {
		fmt.Print("Create message err: ")
		fmt.Println(err)
	}
	
	/* matching check */
	return si.NewPostMessageOK() /* TODO hokanotoissho  {code message}*/
}

func GetMessages(p si.GetMessagesParams) middleware.Responder {
	token := p.Token
	pid := p.UserID
	latest := p.Latest
	oldest := p.Oldest
	limit := *p.Limit

	err := util.ValidateToken(token)
	if err != nil {
		fmt.Println("Invalid token err:")
		fmt.Println(err)
	}

	uid, err := util.GetIDByToken(token)
	
	if err != nil {
		fmt.Print("Get id err: ")
		fmt.Println(err)
	}

	/* TODO validate matching state? */

	r := repositories.NewUserMessageRepository()

	messages, err := r.GetMessages(uid, pid, int(limit), latest, oldest)
	if err != nil {
		fmt.Print("Get messages err: ")
		fmt.Println(err)
	}

	var reses entities.UserMessages
	reses = messages

	return si.NewGetMessagesOK().WithPayload(reses.Build())
}
