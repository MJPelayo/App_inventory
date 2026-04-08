// (placeholder for future)

// Package middleware contains HTTP middleware functions for request processing
package middleware

import (
    "net/http"
)

// CORSMiddleware adds CORS headers to allow frontend requests
// This is needed when frontend runs on different port (e.g., React on 3000, API on 8080)
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// LoggingMiddleware logs each incoming request (useful for debugging)
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Log request (will implement proper logging in final project)
        next.ServeHTTP(w, r)
    })
}