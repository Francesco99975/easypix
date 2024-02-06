package controller

import (
	"net/http"
)

func EnableCORS(fs http.Handler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // You can replace '*' with specific origins
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // If needed for credentials like cookies
		w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin") // If needed for credentials like cookies
        fs.ServeHTTP(w, r)
    }
}
