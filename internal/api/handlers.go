package api

import (
	"encoding/json"
	"net/http"

	"src/internal/services"
)

func GenerateTokensHandler(svc *services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GenerateTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ip, err := getIp(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		accessToken, refreshToken, err := svc.Token.GenerateTokens(req.UserId, ip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func RefreshTokensHandler(svc *services.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data RefreshTokenRequest
		ip, err := getIp(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		accessToken, refreshToken, err := svc.RefreshTokens(data.AccessToken, data.RefreshToken, ip)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}
