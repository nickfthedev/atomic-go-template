package user

import (
	"my-go-template/internal/middleware"
	"my-go-template/internal/model"
	"net/http"
)

func GetUserFromContext(r *http.Request) model.User {
	user, ok := r.Context().Value(middleware.UserKey).(model.User)
	if !ok {
		return model.User{}
	}
	return user
}
