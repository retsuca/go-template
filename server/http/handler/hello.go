package handler

import (
	"encoding/json"
	"net/http"

	"go-template/server/http/types"
)

// @Router       / [get]
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	h.Metrics.HelloCounter.WithLabelValues("test").Inc()
	h.Metrics.HelloGauge.WithLabelValues("test").Set(1)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := types.HelloResponse{
		Message: "Hello, World!",
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// @Param name query string true "name"
func (h *Handler) HelloWithParam(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := types.HelloWithParamResponse{
		Message: name,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
