package middleware

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// Recoverer is a middleware that recovers from panics and returns a 500 error
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				zap.L().Error("Panic recovered", 
					zap.Any("panic", rvr),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				
				errResp := ErrorResponse{
					Error:   "Internal Server Error",
					Message: "An unexpected error occurred",
					Code:    http.StatusInternalServerError,
				}

				if err := json.NewEncoder(w).Encode(errResp); err != nil {
					zap.L().Error("Failed to write error response", zap.Error(err))
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
} 