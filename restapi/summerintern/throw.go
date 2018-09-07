package summerintern

const debugMode = false

func debugMessage(err error) string {
	if debugMode && err != nil {
		return ": " + err.Error()
	}
	return ""
}

func PostMessageThrowInternalServerError(err error) *PostMessageInternalServerError {
	return NewPostMessageInternalServerError().WithPayload(
		&PostMessageInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error" + debugMessage(err),
		})
}

func PostMessageThrowUnauthorized() *PostMessageUnauthorized {
	return NewPostMessageUnauthorized().WithPayload(
		&PostMessageUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗)",
		})
}

func PostMessageThrowBadRequest(mes string) *PostMessageBadRequest {
	return NewPostMessageBadRequest().WithPayload(
		&PostMessageBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetMessagesThrowInternalServerError(err error) *GetMessagesInternalServerError {
	return NewGetMessagesInternalServerError().WithPayload(
		&GetMessagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error" + debugMessage(err),
		})
}

func GetMessagesThrowUnauthorized() *GetMessagesUnauthorized {
	return NewGetMessagesUnauthorized().WithPayload(
		&GetMessagesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗)",
		})
}

func GetMessagesThrowBadRequest(mes string) *GetMessagesBadRequest {
	return NewGetMessagesBadRequest().WithPayload(
		&GetMessagesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetTokenByUserIDThrowInternalServerError(err error) *GetTokenByUserIDInternalServerError {
	return NewGetTokenByUserIDInternalServerError().WithPayload(
		&GetTokenByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error" + debugMessage(err),
		})
}

func GetTokenByUserIDThrowNotFound() *GetTokenByUserIDNotFound {
	return NewGetTokenByUserIDNotFound().WithPayload(
		&GetTokenByUserIDNotFoundBody{
			Code:    "404",
			Message: "User Token Not Found",
		})
}

func GetUsersThrowInternalServerError(err error) *GetUsersInternalServerError {
	return NewGetUsersInternalServerError().WithPayload(
		&GetUsersInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error" + debugMessage(err),
		})
}

func GetUsersThrowUnauthorized() *GetUsersUnauthorized {
	return NewGetUsersUnauthorized().WithPayload(
		&GetUsersUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗)",
		})
}

func GetUsersThrowBadRequest(mes string) *GetUsersBadRequest {
	return NewGetUsersBadRequest().WithPayload(
		&GetUsersBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetProfileByUserIDThrowInternalServerError(err error) *GetProfileByUserIDInternalServerError {
	return NewGetProfileByUserIDInternalServerError().WithPayload(
		&GetProfileByUserIDInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error" + debugMessage(err),
		})
}

func GetProfileByUserIDThrowUnauthorized() *GetProfileByUserIDUnauthorized {
	return NewGetProfileByUserIDUnauthorized().WithPayload(
		&GetProfileByUserIDUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗)",
		})
}

func GetProfileByUserIDThrowBadRequest(mes string) *GetProfileByUserIDBadRequest {
	return NewGetProfileByUserIDBadRequest().WithPayload(
		&GetProfileByUserIDBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetProfileByUserIDThrowNotFound() *GetProfileByUserIDNotFound {
	return NewGetProfileByUserIDNotFound().WithPayload(
		&GetProfileByUserIDNotFoundBody{
			Code:    "400",
			Message: "User Not Found. (そのIDのユーザーは存在しません.)",
		})
}

func PutProfileThrowInternalServerError(err error) *PutProfileInternalServerError {
	return NewPutProfileInternalServerError().WithPayload(
		&PutProfileInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error" + debugMessage(err),
		})
}

func PutProfileThrowUnauthorized() *PutProfileUnauthorized {
	return NewPutProfileUnauthorized().WithPayload(
		&PutProfileUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗)",
		})
}

func PutProfileThrowBadRequest(mes string) *PutProfileBadRequest {
	return NewPutProfileBadRequest().WithPayload(
		&PutProfileBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func PutProfileThrowForbidden() *PutProfileForbidden {
	return NewPutProfileForbidden().WithPayload(
		&PutProfileForbiddenBody{
			Code:    "403",
			Message: "Forbidden. (他の人のプロフィールは更新できません.)",
		})
}

func PostImageThrowInternalServerError(err error) *PostImagesInternalServerError {
	return NewPostImagesInternalServerError().WithPayload(
		&PostImagesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error" + debugMessage(err),
		})
}

func PostImageThrowUnauthorized() *PostImagesUnauthorized {
	return NewPostImagesUnauthorized().WithPayload(
		&PostImagesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗)",
		})
}

func PostImageThrowBadRequest(mes string) *PostImagesBadRequest {
	return NewPostImagesBadRequest().WithPayload(
		&PostImagesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetLikesThrowInternalServerError(err error) *GetLikesInternalServerError {
	return NewGetLikesInternalServerError().WithPayload(
		&GetLikesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error" + debugMessage(err),
		})
}

func GetLikesThrowUnauthorized() *GetLikesUnauthorized {
	return NewGetLikesUnauthorized().WithPayload(
		&GetLikesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗)",
		})
}

func GetLikesThrowBadRequest(mes string) *GetLikesBadRequest {
	return NewGetLikesBadRequest().WithPayload(
		&GetLikesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func PostLikeThrowInternalServerError(err error) *PostLikeInternalServerError {
	return NewPostLikeInternalServerError().WithPayload(
		&PostLikeInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error" + debugMessage(err),
		})
}

func PostLikeThrowUnauthorized() *PostLikeUnauthorized {
	return NewPostLikeUnauthorized().WithPayload(
		&PostLikeUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗)",
		})
}

func PostLikeThrowBadRequest(mes string) *PostLikeBadRequest {
	return NewPostLikeBadRequest().WithPayload(
		&PostLikeBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}

func GetMatchesThrowInternalServerError(err error) *GetMatchesInternalServerError {
	return NewGetMatchesInternalServerError().WithPayload(
		&GetMatchesInternalServerErrorBody{
			Code:    "500",
			Message: "Internal Server Error" + debugMessage(err),
		})
}

func GetMatchesThrowUnauthorized() *GetMatchesUnauthorized {
	return NewGetMatchesUnauthorized().WithPayload(
		&GetMatchesUnauthorizedBody{
			Code:    "401",
			Message: "Unauthorized (トークン認証に失敗)",
		})
}

func GetMatchesThrowBadRequest(mes string) *GetMatchesBadRequest {
	return NewGetMatchesBadRequest().WithPayload(
		&GetMatchesBadRequestBody{
			Code:    "400",
			Message: "Bad Request: " + mes,
		})
}
