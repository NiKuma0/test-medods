package services

import (
	"fmt"
	"jwt-service/internal/repositories"
)

type IMailClient interface {
	SendNotification(msg, email string) error
}

type NotificationService struct {
	client IMailClient
	repos  *repositories.Repositories
}

func NewNotificationService(client IMailClient, repos *repositories.Repositories) *NotificationService {
	return &NotificationService{
		client: client,
		repos:  repos,
	}
}

func (s *NotificationService) NewIpEnterNotification(userId, ip string) error {
	user, err := s.repos.User.Get(userId)
	if err != nil {
		return err
	}
	return s.client.SendNotification(
		fmt.Sprintf(
			"You have entered from new IP address (%s)",
			ip,
		),
		user.Email,
	)
}
