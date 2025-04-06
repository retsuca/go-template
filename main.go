/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
// Package main is the entry point for the go-template application.
// It provides a flexible template for building HTTP and gRPC services with
// database support and configuration management.
package main

import (
	"go-template/cmd"
	_ "go-template/internal/clients/httpClient"

	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd.Execute()
}
