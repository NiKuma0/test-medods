package services

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"src/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService struct {
	secret       string
	repos        *repositories.Repositories
	notification INotificationService
}

type INotificationService interface {
	NewIpEnterNotification(userId, ip string) (err error)
}

type TokenRepository interface {
	SaveRefreshToken(userId, tokenHash string) error
	DeleteRefreshToken(tokenHash string) error
	IsRefreshTokenValid(userId, tokenHash string) (bool, error)
}

type Claims struct {
	jwt.RegisteredClaims
	UserId string `json:"user_id"`
	Ip     string `json:"ip"`
	Hash   string `json:"hash"`
}

func NewTokenService(secret string, notificationService INotificationService, repos *repositories.Repositories) TokenService {
	return TokenService{
		secret: secret,
		repos:  repos,
	}
}

func (s *TokenService) generateAccessToken(userId, refreshTokenHash, ip string) (string, error) {
	claims := &Claims{
		UserId: userId,
		Ip:     ip,
		Hash:   refreshTokenHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *TokenService) generateRefreshToken() string {
	return base64.URLEncoding.EncodeToString([]byte(uuid.NewString()))
}

func (s *TokenService) hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func (s *TokenService) ValidateAccessToken(accessToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid access token")
}

func (s *TokenService) GenerateTokens(userId, ip string) (string, string, error) {
	refreshToken := s.generateRefreshToken()
	refreshTokenHash := s.hashRefreshToken(refreshToken)
	err := s.repos.Token.SaveRefreshToken(userId, refreshTokenHash)
	if err != nil {
		return "", "", err
	}
	accessToken, err := s.generateAccessToken(userId, refreshTokenHash, ip)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, err
}

func (s *TokenService) RefreshTokens(accessToken, refreshToken, ip string) (string, string, error) {
	claims, err := s.ValidateAccessToken(accessToken)
	if err != nil {
		return "", "", err
	}
	if claims.Ip != ip {
		defer s.notification.NewIpEnterNotification(claims.UserId, ip)
	}
	refreshTokenHash := s.hashRefreshToken(refreshToken)
	valid, err := s.repos.Token.IsRefreshTokenValid(claims.UserId, refreshTokenHash)
	if err != nil {
		return "", "", err
	}
	if !valid || claims.Hash != refreshTokenHash {
		return "", "", errors.New("refresh token is invalid")
	}
	defer s.repos.Token.DeleteRefreshToken(refreshTokenHash)
	return s.GenerateTokens(claims.UserId, ip)
}
