package routes

import (
	"encoding/json"
	"net/http"

	"go-template/server/http/handler"
	"go-template/server/http/types"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	swagger "go-template/proto/gen/swagger"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// HelloResponse represents the hello endpoint response
type HelloResponse struct {
	Message string `json:"message"`
}

// HelloWithParamResponse represents the hello with param endpoint response
type HelloWithParamResponse struct {
	Message string `json:"message"`
}

// SetupRoutes configures all routes for the server
func SetupRoutes(r *chi.Mux, h *handler.Handler, gwMux *runtime.ServeMux) {
	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// Configure Swagger UI only when gRPC gateway is enabled
	if gwMux != nil {
		// Serve gRPC-Gateway Swagger documentation
		r.Get("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(swagger.ApidocsSwaggerJson)
		})

		r.Handle("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger/swagger.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		))
		r.Handle("/*", gwMux)
	} else {
		r.Handle("/swagger/*", httpSwagger.WrapHandler)
	}

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := types.HealthResponse{
			Status: "healthy",
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	// Application routes
	r.Get("/", h.Hello)
	r.Get("/withparam", h.HelloWithParam)
}
