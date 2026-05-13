package main

import (
	"log"

	"net/http"

	"os"

	"strings"

	httpPkg "tech-ip-sem2/services/tasks/internal/http"

	"tech-ip-sem2/shared/middleware"
)

func main() {

	port := os.Getenv("TASKS_PORT")

	if port == "" {

		port = "8082"

	}

	authURL := os.Getenv("AUTH_BASE_URL")

	if authURL == "" {

		authURL = "http://localhost:8081"

	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/tasks", httpPkg.AuthMiddleware(authURL, func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {

			httpPkg.CreateTaskHandler(w, r)

		} else if r.Method == "GET" {

			httpPkg.ListTasksHandler(w, r)

		} else {

			http.Error(w, "Method not allowed", 405)

		}

	}))

	mux.HandleFunc("/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {

		id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")

		if id == "" {

			http.NotFound(w, r)

			return

		}

		httpPkg.AuthMiddleware(authURL, func(w http.ResponseWriter, r *http.Request) {

			switch r.Method {

			case "GET":

				httpPkg.GetTaskHandler(w, r, id)

			case "PATCH":

				httpPkg.UpdateTaskHandler(w, r, id)

			case "DELETE":

				httpPkg.DeleteTaskHandler(w, r, id)

			default:

				http.Error(w, "Method not allowed", 405)

			}

		})(w, r)

	})

	handler := middleware.LoggingMiddleware(middleware.RequestIDMiddleware(mux))

	log.Printf("Tasks service starting on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, handler))

}
