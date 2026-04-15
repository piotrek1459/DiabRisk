//go:build integration

package main

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestDatabaseIntegration_MigrationsCreateExpectedSchema(t *testing.T) {
	dbURL := os.Getenv("DATA_INTEGRATION_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		t.Skip("set DATA_INTEGRATION_DATABASE_URL or DATABASE_URL to run data-svc integration tests")
	}

	if err := waitForDB(dbURL); err != nil {
		t.Fatalf("database did not become ready: %v", err)
	}
	if err := runMigrations(dbURL); err != nil {
		t.Fatalf("migrations failed: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("failed to create database pool: %v", err)
	}
	defer pool.Close()

	if err := verifySchema(pool); err != nil {
		t.Fatalf("schema verification failed: %v", err)
	}
}
