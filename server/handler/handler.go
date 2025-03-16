package handler

import (
	"context"
)

type Client interface {
	Do(ctx context.Context, method, path string, body any, args map[string]string) ([]byte, error)
}

type Handler struct {
	HTTPClient Client
}

func NewHandler(client Client) *Handler {
	return &Handler{
		HTTPClient: client,
	}
}
