package userimage

import (
	si "github.com/eure/si2018-server-side/restapi/summerintern"
	"github.com/go-openapi/runtime/middleware"
)

type Validator interface {
	Validate() middleware.Responder
}

type PostValidator struct {
	token       string
	imageBase64 []byte
}

func NewPostValidator(t string, i []byte) Validator {
	return PostValidator{
		token:       t,
		imageBase64: i,
	}
}

func (v PostValidator) Validate() middleware.Responder {
	if len(v.imageBase64) == 0 {
		return si.NewPostImagesBadRequest().WithPayload(
			&si.PostImagesBadRequestBody{
				Code:    "400",
				Message: "Bad Request",
			})
	}

	if len(v.token) == 0 {
		return si.NewPostImagesUnauthorized().WithPayload(
			&si.PostImagesUnauthorizedBody{
				Code:    "401",
				Message: "Unauthorized",
			})
	}

	return nil
}
