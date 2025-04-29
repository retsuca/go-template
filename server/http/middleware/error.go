package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// ErrorHandler is a middleware that handles errors and returns a standardized error response
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture the status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Call the next handler
		next.ServeHTTP(ww, r)

		// Check if we have an error status code
		if ww.Status() >= 400 {
			// Create error response
			errResp := ErrorResponse{
				Error: http.StatusText(ww.Status()),
				Code:  ww.Status(),
			}

			// Set content type
			w.Header().Set("Content-Type", "application/json")

			// Write the error response
			if err := json.NewEncoder(w).Encode(errResp); err != nil {
				zap.L().Error("Failed to write error response", zap.Error(err))
			}
		}
	})
}
