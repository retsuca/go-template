package types

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