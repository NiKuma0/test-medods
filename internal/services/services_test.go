package services_test

import (
	"errors"
	"testing"

	"src/internal/repositories"
	"src/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMailClient struct {
	mock.Mock
}

func (m *MockMailClient) SendNotification(msg, email string) error {
	args := m.Called(msg, email)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Get(userId string) (repositories.User, error) {
	args := m.Called(userId)
	return args.Get(0).(repositories.User), args.Error(1)
}

func (m *MockUserRepository) IsExists(userId string) (bool, error) {
	args := m.Called(userId)
	return args.Get(0).(bool), args.Error(1)
}

func TestNotificationService_NewIpEnterNotification(t *testing.T) {
	mockMailClient := new(MockMailClient)
	mockUserRepo := new(MockUserRepository)

	notificationService := services.NewNotificationService(mockMailClient, &repositories.Repositories{
		User: mockUserRepo,
	})

	t.Run("Success", func(t *testing.T) {
		userId := "123e4567-e89b-12d3-a456-426614174000"
		ip := "192.168.1.100"
		user := repositories.User{
			UserId: userId,
			Email:  "test@example.com",
		}
		mockUserRepo.On("Get", userId).Return(user, nil)
		mockMailClient.On("SendNotification", mock.Anything, user.Email).Return(nil)

		err := notificationService.NewIpEnterNotification(userId, ip)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockMailClient.AssertExpectations(t)
	})
}

// Mock for TokenRepository
type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) SaveRefreshToken(userId, tokenHash string) error {
	args := m.Called(userId, tokenHash)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteRefreshToken(tokenHash string) error {
	args := m.Called(tokenHash)
	return args.Error(0)
}

func (m *MockTokenRepository) IsRefreshTokenValid(userId, tokenHash string) (bool, error) {
	args := m.Called(userId, tokenHash)
	return args.Bool(0), args.Error(1)
}

// Mock for INotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) NewIpEnterNotification(userId, ip string) error {
	args := m.Called(userId, ip)
	return args.Error(0)
}

func TestGenerateAccessToken(t *testing.T) {
	mockRepo := new(MockTokenRepository)
	mockNotService := new(MockNotificationService)
	repos := repositories.Repositories{
		Token: mockRepo,
	}

	service := services.NewTokenService(
		"secret",
		mockNotService,
		&repos,
	)

	t.Run("GenerateTokens", func(t *testing.T) {
		userId := "someID"
		ip := "127.0.0.1"

		mockRepo.On("SaveRefreshToken", userId, mock.Anything).Return(nil)
		accessToken, refreshToken, err := service.GenerateTokens(userId, ip)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		assert.NotEqual(t, "", accessToken)
		assert.NotEqual(t, "", refreshToken)
	})

	t.Run("GenerateTokensSaveNotSucceed", func(t *testing.T) {
		userId := "someID"
		ip := "127.0.0.1"

		mockRepo.On("SaveRefreshToken", userId, mock.Anything).Return(errors.New("something went wrong"))
		_, _, err := service.GenerateTokens(userId, ip)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RefreshTokens", func(t *testing.T) {
		userId := "someID"
		ip := "127.0.0.1"

		mockRepo.On("SaveRefreshToken", userId, mock.Anything).Return(nil)
		accessToken, refreshToken, _ := service.GenerateTokens(userId, ip)

		mockNotService.On("NewIpEnterNotification", userId, ip).Return(nil)
		mockRepo.On("SaveRefreshToken", userId, mock.Anything).Return(nil)
		mockRepo.On("IsRefreshTokenValid", userId, mock.Anything).Return(true, nil)
		mockRepo.On("DeleteRefreshToken", mock.Anything).Return(nil)
		newAccessToken, newRefreshToken, err := service.RefreshTokens(accessToken, refreshToken, ip)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockNotService.AssertNotCalled(t, "NewIpEnterNotification")
		assert.NotEqual(t, "", newAccessToken)
		assert.NotEqual(t, "", newRefreshToken)
	})

	t.Run("RefreshTokens_NotValidTokens", func(t *testing.T) {
		userId := "someID"
		ip := "127.0.0.1"

		mockRepo.On("SaveRefreshToken", userId, mock.Anything).Return(nil)
		accessToken, refreshToken := "NotValidAccessToken", "NotValidRefreshToken"

		mockNotService.On("NewIpEnterNotification", userId, ip).Return(nil)
		mockRepo.On("SaveRefreshToken", userId, mock.Anything).Return(nil)
		mockRepo.On("IsRefreshTokenValid", userId, mock.Anything).Return(true, nil)
		mockRepo.On("DeleteRefreshToken", mock.Anything).Return(nil)
		newAccessToken, newRefreshToken, err := service.RefreshTokens(accessToken, refreshToken, ip)
		assert.Error(t, err, "invalid access token")
		mockRepo.AssertNotCalled(t, "SaveRefreshToken")
		mockRepo.AssertNotCalled(t, "IsRefreshTokenValid")
		mockRepo.AssertNotCalled(t, "DeleteRefreshToken")
		mockNotService.AssertNotCalled(t, "NewIpEnterNotification")
		assert.Equal(t, "", newAccessToken)
		assert.Equal(t, "", newRefreshToken)
	})
}
