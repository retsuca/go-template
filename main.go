package main

import (
	logger "go-template/pkg/logger"

	"go-template/internal/config"
	"go-template/server"
)

func main() {

	defer logger.Sync() // flushes buffer, if any

	server.CreateHTPPServer(config.Get(config.HTTP_HOST), config.Get(config.HTTP_PORT))
}
