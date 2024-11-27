package main

import (
	"log"
	"net/http"

	"jwt-service/internal/application"
)

func main() {
	app := application.NewApplication()
	defer app.Shutdown()
	router := app.Controller.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
