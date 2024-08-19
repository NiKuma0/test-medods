package application

import (
	"fmt"
	"src/internal/clients"
	"src/internal/repositories"
	"src/internal/services"
)

type Application struct {
	Services      *services.Services
	Repositories  *repositories.Repositories
	Config        *Config
	shutdownStack []func()
}

func NewApplication() *Application {
	config := NewConfig()
	userRepo := repositories.NewUserRepositoryFromDataSource(config.POSTGRES_DSN)
	tokenRepo := repositories.NewRefreshTokenRepositoryFromDataSource(config.POSTGRES_DSN)
	repos := repositories.Repositories{
		User:  &userRepo,
		Token: &tokenRepo,
	}

	smtpMockServer := clients.StartNewMockServer()
	mailClient, err := clients.NewMailClientNoAuth(
		config.SMTP_SENDER_USERNAME,
		fmt.Sprintf("%s:%d", "127.0.0.1", smtpMockServer.PortNumber()),
	)
	if err != nil {
		panic(err)
	}

	userService := services.NewUserService(&repos)
	notificationService := services.NewNotificationService(&mailClient, &repos)
	tokenService := services.NewTokenService("secret", &notificationService, &repos)

	return &Application{
		Config: &config,
		Services: &services.Services{
			Notification: &notificationService,
			Token:        &tokenService,
			User:         &userService,
		},
		shutdownStack: []func(){
			userRepo.Shutdown,
			mailClient.Shutdown,
			func() { smtpMockServer.Stop() },
		},
	}
}

func (a *Application) Shutdown() {
	for _, f := range a.shutdownStack {
		defer f()
	}
}
