/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
// Package cmd provides command-line interface functionality for the application.
package cmd

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Required for file-based migrations
	"github.com/spf13/cobra"
	"go-template/internal/clients/db"
	"go-template/pkg/logger"
	"go.uber.org/zap"
)

const (
	migrationsPath = "file://database/migrations"
	dbDriver       = "postgres"
)

// DatabaseMigrationCmd represents the root database-migration command.
var DatabaseMigrationCmd = &cobra.Command{
	Use:   "database-migration",
	Short: "Manage database migrations",
	Long: `Database migration command provides functionality to manage database schema migrations.
It supports both applying (up) and reverting (down) migrations.`,
}

// DatabaseMigrationUpCmd represents the command to apply migrations.
var DatabaseMigrationUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply database migrations",
	Long:  `Apply all pending database migrations to update the schema to the latest version.`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := databaseMigrationUp(); err != nil {
			logger.Fatal("Failed to apply migrations", zap.Error(err))
		}
		logger.Info("Successfully applied all migrations")
	},
}

// databaseMigrationUp handles the database migration up operation.
func databaseMigrationUp() error {
	m, err := initializeMigration()
	if err != nil {
		return fmt.Errorf("failed to initialize migration: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No migrations to apply - database is up to date")

			return nil
		}

		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// DatabaseMigrationDownCmd represents the command to revert migrations.
var DatabaseMigrationDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Revert database migrations",
	Long:  `Revert all applied database migrations to downgrade the schema to its base version.`,
	Run: func(_ *cobra.Command, _ []string) {
		if err := databaseMigrationDown(); err != nil {
			logger.Fatal("Failed to revert migrations", zap.Error(err))
		}
		logger.Info("Successfully reverted all migrations")
	},
}

// databaseMigrationDown handles the database migration down operation.
func databaseMigrationDown() error {
	m, err := initializeMigration()
	if err != nil {
		return fmt.Errorf("failed to initialize migration: %w", err)
	}

	if err := m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No migrations to revert - database is at base version")

			return nil
		}

		return fmt.Errorf("failed to revert migrations: %w", err)
	}

	return nil
}

// initializeMigration creates and configures a new migration instance.
func initializeMigration() (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, dbDriver, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return m, nil
}

func init() {
	DatabaseMigrationCmd.AddCommand(DatabaseMigrationUpCmd)
	DatabaseMigrationCmd.AddCommand(DatabaseMigrationDownCmd)
	rootCmd.AddCommand(DatabaseMigrationCmd)
}
