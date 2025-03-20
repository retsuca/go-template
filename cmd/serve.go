/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"go-template/internal/config"
	"go-template/pkg/logger"
	serverGRPC "go-template/server/grpc"
	serverHTTP "go-template/server/http"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "create a server",
}

var serveHTTPCmd = &cobra.Command{
	Use:   "http",
	Short: "Starts an echo API Server",
	Run: func(cmd *cobra.Command, _ []string) {
		serveHTTP(cmd.Context())

		fmt.Println("serve called")
	},
}

func serveHTTP(ctx context.Context) {
	defer logger.Sync() // flushes buffer, if any
	serverHTTP.CreateHTPPServer(ctx, config.Get(config.HOST), config.Get(config.HTTP_PORT))
}

// serveCmd represents the serve command.
var serveGRPCCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Starts a grpc Server",
	Run: func(cmd *cobra.Command, _ []string) {
		serveGRPC(cmd.Context())

		fmt.Println("serve called")
	},
}

func serveGRPC(ctx context.Context) {
	defer logger.Sync() // flushes buffer, if any
	serverGRPC.CreateGRPCServer(ctx, config.Get(config.HOST), config.Get(config.GRPC_PORT), config.Get(config.HTTP_PORT))
}

func init() {
	serveCmd.AddCommand(serveHTTPCmd)
	serveCmd.AddCommand(serveGRPCCmd)
	rootCmd.AddCommand(serveCmd)
}
