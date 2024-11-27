package api

import (
	"github.com/gorilla/mux"
)

func (c *Controller) NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/token", c.GenerateTokensHandler).Methods("POST")
	router.HandleFunc("/refresh", c.RefreshTokensHandler).Methods("POST")

	router.Use(loggingMiddleware)
	router.Use(enableCors)
	router.Use(contentTypeApplicationJsonMiddleware)

	return router
}
