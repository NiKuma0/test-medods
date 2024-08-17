package main

import (
	"log"
	"net/http"

	"src/internal/api"
	"src/internal/application"
)

func main() {
	app := application.NewApplication()
	defer app.Shutdown()
	router := api.NewRouter(app)
	log.Fatal(http.ListenAndServe(":8080", router))
}
