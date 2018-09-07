package summerintern

func PostMessageThrowInternalServerError(fun string, err error) *PostMessageInternalServerError {
	return NewPostMessageInternalServerError().WithPayload(
		&PostMessageInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func PostMessageThrowUnauthorized(mes string) *PostMessageUnauthorized {
	return NewPostMessageUnauthorized().WithPayload(
		&PostMessageUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func PostMessageThrowBadRequest(mes string) *PostMessageBadRequest {
	return NewPostMessageBadRequest().WithPayload(
		&PostMessageBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetMessagesThrowInternalServerError(fun string, err error) *GetMessagesInternalServerError {
	return NewGetMessagesInternalServerError().WithPayload(
		&GetMessagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func GetMessagesThrowUnauthorized(mes string) *GetMessagesUnauthorized {
	return NewGetMessagesUnauthorized().WithPayload(
		&GetMessagesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func GetMessagesThrowBadRequest(mes string) *GetMessagesBadRequest {
	return NewGetMessagesBadRequest().WithPayload(
		&GetMessagesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetTokenByUserIDThrowInternalServerError(fun string, err error) *GetTokenByUserIDInternalServerError {
	return NewGetTokenByUserIDInternalServerError().WithPayload(
		&GetTokenByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func GetTokenByUserIDThrowNotFound(mes string) *GetTokenByUserIDNotFound {
	return NewGetTokenByUserIDNotFound().WithPayload(
		&GetTokenByUserIDNotFoundBody{
			Code:    "404",
			Message: "User Token Not Found: " + mes,
		})
}

func GetUsersThrowInternalServerError(fun string, err error) *GetUsersInternalServerError {
	return NewGetUsersInternalServerError().WithPayload(
		&GetUsersInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func GetUsersThrowUnauthorized(mes string) *GetUsersUnauthorized {
	return NewGetUsersUnauthorized().WithPayload(
		&GetUsersUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func GetUsersThrowBadRequest(mes string) *GetUsersBadRequest {
	return NewGetUsersBadRequest().WithPayload(
		&GetUsersBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetProfileByUserIDThrowInternalServerError(fun string, err error) *GetProfileByUserIDInternalServerError {
	return NewGetProfileByUserIDInternalServerError().WithPayload(
		&GetProfileByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func GetProfileByUserIDThrowUnauthorized(mes string) *GetProfileByUserIDUnauthorized {
	return NewGetProfileByUserIDUnauthorized().WithPayload(
		&GetProfileByUserIDUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func GetProfileByUserIDThrowBadRequest(mes string) *GetProfileByUserIDBadRequest {
	return NewGetProfileByUserIDBadRequest().WithPayload(
		&GetProfileByUserIDBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetProfileByUserIDThrowNotFound(mes string) *GetProfileByUserIDNotFound {
	return NewGetProfileByUserIDNotFound().WithPayload(
		&GetProfileByUserIDNotFoundBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func PutProfileThrowInternalServerError(fun string, err error) *PutProfileInternalServerError {
	return NewPutProfileInternalServerError().WithPayload(
		&PutProfileInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func PutProfileThrowUnauthorized(mes string) *PutProfileUnauthorized {
	return NewPutProfileUnauthorized().WithPayload(
		&PutProfileUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func PutProfileThrowBadRequest(mes string) *PutProfileBadRequest {
	return NewPutProfileBadRequest().WithPayload(
		&PutProfileBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func PutProfileThrowForbidden(mes string) *PutProfileForbidden {
	return NewPutProfileForbidden().WithPayload(
		&PutProfileForbiddenBody{
			Code:    "403",
			Message: "Forbidden. (他の人のプロフィールは更新できません.): " + mes,
		})
}

func PostImageThrowInternalServerError(fun string, err error) *PostImagesInternalServerError {
	return NewPostImagesInternalServerError().WithPayload(
		&PostImagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func PostImageThrowUnauthorized(mes string) *PostImagesUnauthorized {
	return NewPostImagesUnauthorized().WithPayload(
		&PostImagesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func PostImageThrowBadRequest(mes string) *PostImagesBadRequest {
	return NewPostImagesBadRequest().WithPayload(
		&PostImagesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetLikesThrowInternalServerError(fun string, err error) *GetLikesInternalServerError {
	return NewGetLikesInternalServerError().WithPayload(
		&GetLikesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func GetLikesThrowUnauthorized(mes string) *GetLikesUnauthorized {
	return NewGetLikesUnauthorized().WithPayload(
		&GetLikesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func GetLikesThrowBadRequest(mes string) *GetLikesBadRequest {
	return NewGetLikesBadRequest().WithPayload(
		&GetLikesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func PostLikeThrowInternalServerError(fun string, err error) *PostLikeInternalServerError {
	return NewPostLikeInternalServerError().WithPayload(
		&PostLikeInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func PostLikeThrowUnauthorized(mes string) *PostLikeUnauthorized {
	return NewPostLikeUnauthorized().WithPayload(
		&PostLikeUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func PostLikeThrowBadRequest(mes string) *PostLikeBadRequest {
	return NewPostLikeBadRequest().WithPayload(
		&PostLikeBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetMatchesThrowInternalServerError(fun string, err error) *GetMatchesInternalServerError {
	return NewGetMatchesInternalServerError().WithPayload(
		&GetMatchesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error: " + fun + " failed: " + err.Error(),
		})
}

func GetMatchesThrowUnauthorized(mes string) *GetMatchesUnauthorized {
	return NewGetMatchesUnauthorized().WithPayload(
		&GetMatchesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗): " + mes,
		})
}

func GetMatchesThrowBadRequest(mes string) *GetMatchesBadRequest {
	return NewGetMatchesBadRequest().WithPayload(
		&GetMatchesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}
