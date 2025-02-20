/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"go-template/internal/config"
	"go-template/pkg/logger"
	"go-template/server"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts an echo API Server",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())

		fmt.Println("serve called")
	},
}

func serve(ctx context.Context) {
	defer logger.Sync() // flushes buffer, if any
	server.CreateHTPPServer(config.Get(config.HTTP_HOST), config.Get(config.HTTP_PORT))

}
func init() {
	rootCmd.AddCommand(serveCmd)
g}
