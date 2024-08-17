package services

import (
	"fmt"
	"src/internal/clients"
	"src/internal/repositories"
)

type NotificationService struct {
	client *clients.MailClient
	repos  *repositories.Repositories
}

func NewNotificationService(client *clients.MailClient, repos *repositories.Repositories) NotificationService {
	return NotificationService{
		client: client,
		repos:  repos,
	}
}

func (s *NotificationService) NewIpEnterNotification(userId, ip string) (err error) {
	user, err := s.repos.User.Get(userId)
	if err != nil {
		return
	}
	return s.client.SendNotification(
		fmt.Sprintf(
			"You have entered from new IP address (%s)",
			ip,
		),
		user.Email,
	)
}
