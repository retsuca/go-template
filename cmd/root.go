/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
// Package cmd provides command-line interface functionality for the application.
package cmd

import (
	"os"

	"go-template/pkg/logger"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "go-template",
	Short: "A modern Go service template with HTTP and gRPC support",
	Long: `go-template is a production-ready service template that provides:

- HTTP and gRPC server support with automatic OpenAPI documentation
- Database migrations and management
- Structured logging and metrics
- Configuration management
- Graceful shutdown handling
- Health checks and monitoring endpoints`,
	Version: "1.0.0",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Failed to execute root command", zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	// Configure persistent flags that will be inherited by all subcommands
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file path (default is .env)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug mode")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose logging")

	// Bind flags to environment variables
	if err := rootCmd.PersistentFlags().MarkHidden("debug"); err != nil {
		logger.Error("Failed to mark debug flag as hidden", zap.Error(err))
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
