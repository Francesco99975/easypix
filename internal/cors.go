package controller

import "net/http"

func EnableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "https://francescobarranca.dev,https://www.francescobarranca.dev,http://localhost:3000") // You can replace '*' with specific origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // If needed for credentials like cookies

		// Call the original handler
		handler.ServeHTTP(w, r)
	})
}