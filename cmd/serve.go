/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
// Package cmd provides the command-line interface for the application.
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"go-template/internal/config"
	"go-template/pkg/logger"
	serverGRPC "go-template/server/grpc"
	serverHTTP "go-template/server/http"
)

// serveCmd represents the base serve command.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP or gRPC server",
	Long:  `Serve command allows you to start either an HTTP API server or a gRPC server.`,
}

// serveHTTPCmd represents the HTTP server command.
var serveHTTPCmd = &cobra.Command{
	Use:   "http",
	Short: "Start the HTTP API Server",
	Long:  `Start an HTTP API Server with the configured host and port from environment variables.`,
	Run: func(cmd *cobra.Command, _ []string) {
		serveHTTP(cmd.Context())
	},
}

// serveHTTP initializes and starts the HTTP server.
func serveHTTP(ctx context.Context) {
	defer logger.Sync() // flushes buffer, if any

	host := config.Get(config.HOST)
	port := config.Get(config.HTTP_PORT)

	serverHTTP.CreateHTPPServer(ctx, host, port, nil)
}

// serveGRPCCmd represents the gRPC server command.
var serveGRPCCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Start the gRPC Server",
	Long:  `Start a gRPC Server with the configured host and ports from environment variables.`,
	Run: func(cmd *cobra.Command, _ []string) {
		serveGRPC(cmd.Context())
	},
}

// serveGRPC initializes and starts the gRPC server.
func serveGRPC(ctx context.Context) {
	defer logger.Sync() // flushes buffer, if any

	host := config.Get(config.HOST)
	grpcPort := config.Get(config.GRPC_PORT)
	httpPort := config.Get(config.HTTP_PORT)

	serverGRPC.CreateGRPCServer(ctx, host, grpcPort, httpPort)
}

func init() {
	serveCmd.AddCommand(serveHTTPCmd)
	serveCmd.AddCommand(serveGRPCCmd)
	rootCmd.AddCommand(serveCmd)
}
