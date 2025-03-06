/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"go-template/cmd"
	_ "go-template/internal/clients/http"
	_ "go-template/pkg/metrics"
)

func main() {
	cmd.Execute()
}
