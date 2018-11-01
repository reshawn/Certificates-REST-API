package main

import (
	"Certificates-REST-API/certificate_logic"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	//port to be used variable
	port := "8080"
	//create routes, initialize endpoints
	router := certificates.NewRouter()

	//CORS Settings
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})

	// Launch with CORS
	http.ListenAndServe(":"+port, handlers.CORS(allowedOrigins, allowedMethods)(router))
}
