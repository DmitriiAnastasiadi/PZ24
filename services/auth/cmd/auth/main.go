package main

import (
	"log"

	"net/http"

	"os"

	httpHandlers "tech-ip-sem2/services/auth/internal/http"

	"tech-ip-sem2/shared/middleware"
)

func main() {

	port := os.Getenv("AUTH_PORT")

	if port == "" {

		port = "8081"

	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/auth/login", httpHandlers.LoginHandler)

	mux.HandleFunc("/v1/auth/verify", httpHandlers.VerifyHandler)

	handler := middleware.LoggingMiddleware(middleware.RequestIDMiddleware(mux))

	log.Printf("Auth service starting on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, handler))

}
