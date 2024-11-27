package api

import "jwt-service/internal/services"

type Controller struct {
	notification *services.NotificationService
	token        *services.TokenService
	user         *services.UserService
}

func NewController(
	notification *services.NotificationService,
	token *services.TokenService,
	user *services.UserService,
) *Controller {
	return &Controller{
		notification,
		token,
		user,
	}
}
