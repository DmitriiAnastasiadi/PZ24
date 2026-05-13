package http

import (
	"io"

	"log"

	"net/http"

	"time"
)

func AuthMiddleware(authURL string, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		client := &http.Client{Timeout: 3 * time.Second}

		req, _ := http.NewRequestWithContext(r.Context(), "GET", authURL+"/v1/auth/verify", nil)

		req.Header.Set("Authorization", r.Header.Get("Authorization"))

		if requestID := r.Context().Value("requestID"); requestID != nil {

			req.Header.Set("X-Request-ID", requestID.(string))

		}

		resp, err := client.Do(req)

		if err != nil {

			log.Printf("Auth verify error: %v", err)

			http.Error(w, "Internal server error", 500)

			return

		}

		defer resp.Body.Close()

		if resp.StatusCode == 200 {

			next(w, r)

		} else {

			w.WriteHeader(resp.StatusCode)

			io.Copy(w, resp.Body)

		}

	}

}
