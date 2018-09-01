package userlike

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

type Validator interface {
	Validate() middleware.Responder
}

type GetValidator struct {
	token  string
	limit  int64
	offset int64
}

func NewGetValidator(t string, l, o int64) Validator {
	return GetValidator{
		token:  t,
		limit:  l,
		offset: o,
	}
}

func (v GetValidator) Validate() middleware.Responder {
	if v.limit == 0 {
		return si.NewGetLikesBadRequest().WithPayload(
			&si.GetLikesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	if len(v.token) == 0 {
		return si.NewGetLikesUnauthorized().WithPayload(
			&si.GetLikesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	return nil
}

type PostValidator struct {
	token     string
	partnerID int64
}

func NewPostValidator(t string, id int64) Validator {
	return PostValidator{
		token:     t,
		partnerID: id,
	}
}

func (v PostValidator) Validate() middleware.Responder {
	if v.partnerID == 0 {
		return si.NewPostLikeBadRequest().WithPayload(
			&si.PostLikeBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	if len(v.token) == 0 {
		return si.NewPostLikeUnauthorized().WithPayload(
			&si.PostLikeUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})

	}
	return nil
}
