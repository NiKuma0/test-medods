package api

import (
	"github.com/gorilla/mux"

	"src/internal/application"
)

func NewRouter(app *application.Application) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/token", GenerateTokensHandler(app.Services)).Methods("POST")
	router.HandleFunc("/refresh", RefreshTokensHandler(app.Services.Token)).Methods("POST")

	router.Use(loggingMiddleware)
	router.Use(enableCors)
	router.Use(contentTypeApplicationJsonMiddleware)

	return router
}
