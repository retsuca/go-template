/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"go-template/internal/clients/db"
	"go-template/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

// database-migrationCmd represents the database-migration command
var DatabaseMigrationCmd = &cobra.Command{
	Use:   "database-migration",
	Short: "Do migrations on database",
}

var DatabaseMigrationUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Do migrations on database",
	Run: func(cmd *cobra.Command, args []string) {
		databaseMigrationUp()
	},
}

func databaseMigrationUp() {

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.FatalErr("error when try to connect DB", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	if err != nil {
		logger.FatalErr("error when try to connect DB", err)
	}

	err = m.Up()

	if err != nil {
		logger.FatalErr("error when try to connect DB", err)
	}
}

var DatabaseMigrationDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Down migrations on database",
	Run: func(cmd *cobra.Command, args []string) {
		databaseMigrationDown()
	},
}

func databaseMigrationDown() {

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.FatalErr("error when try to connect DB", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	if err != nil {
		logger.FatalErr("error when try to connect DB", err)
	}

	err = m.Down()

	if err != nil {
		logger.FatalErr("error when try to connect DB", err)
	}
}

func init() {

	DatabaseMigrationCmd.AddCommand(DatabaseMigrationUpCmd)
	DatabaseMigrationCmd.AddCommand(DatabaseMigrationDownCmd)
	rootCmd.AddCommand(DatabaseMigrationCmd)

}
