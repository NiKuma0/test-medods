package services

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	notification *NotificationService
}

func NewTokenService(notification *NotificationService) TokenService {
	return TokenService{
		notification: notification,
	}
}

func (s *TokenService) GenerateTokens(userId, ip string) (accessToken string, refreshToken string, err error) {
	accessToken, err = s.generateAccessToken(userId, ip)
	if err != nil {
		return
	}

	refreshToken, err = s.generateRefreshToken(accessToken)
	if err != nil {
		return
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) generateAccessToken(userId, ip string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"ip":      ip,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte("your-secret-key"))
}

func (s *TokenService) generateRefreshToken(accessToken string) (string, error) {
	refreshToken := base64.StdEncoding.EncodeToString([]byte(accessToken))
	return refreshToken, nil
}

func (s *TokenService) RefreshTokens(accessToken, refreshToken, ip string) (newAccessToken string, newRefreshToken string, err error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("your-secret-key"), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	userId := claims["user_id"].(string)
	if claims["ip"].(string) != ip {
		defer s.notification.NewIpEnterNotification(userId, ip)
	}
	expectedRefreshToken := base64.StdEncoding.EncodeToString([]byte(userId + accessToken + time.Now().String()))
	if err := bcrypt.CompareHashAndPassword([]byte(expectedRefreshToken), []byte(refreshToken)); err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	newAccessToken, err = s.generateAccessToken(userId, ip)
	if err != nil {
		return
	}

	newRefreshToken, err = s.generateRefreshToken(newAccessToken)
	if err != nil {
		return
	}

	return
}
