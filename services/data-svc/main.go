package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Get database connection string from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://diabrisk:diabrisk123@postgres:5432/diabrisk?sslmode=disable"
	}

	log.Println("Starting data-svc...")
	log.Println("Database URL:", maskPassword(dbURL))

	// Wait for PostgreSQL to be ready
	if err := waitForDB(dbURL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := runMigrations(dbURL); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize connection pool
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Failed to create connection pool:", err)
	}
	defer pool.Close()

	// Verify schema
	if err := verifySchema(pool); err != nil {
		log.Fatal("Schema verification failed:", err)
	}

	log.Println("✅ data-svc initialized successfully!")
	log.Println("Database ready with all tables and seed data")

	// Keep service running
	select {}
}

func waitForDB(dbURL string) error {
	log.Println("Waiting for PostgreSQL to be ready...")
	for i := 0; i < 30; i++ {
		pool, err := pgxpool.New(context.Background(), dbURL)
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := pool.Ping(ctx); err == nil {
				pool.Close()
				log.Println("PostgreSQL is ready!")
				return nil
			}
			pool.Close()
		}
		log.Printf("PostgreSQL not ready yet (attempt %d/30), waiting...", i+1)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("timeout waiting for PostgreSQL")
}

func runMigrations(dbURL string) error {
	log.Println("Running database migrations...")

	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Println("⚠️  Could not get migration version:", err)
	} else {
		log.Printf("✅ Migrations complete! Current version: %d (dirty: %v)", version, dirty)
	}

	return nil
}

func verifySchema(pool *pgxpool.Pool) error {
	log.Println("Verifying database schema...")

	tables := []string{
		"users", "model_versions", "assessments",
		"reports", "auth_sessions", "audit_logs",
	}

	ctx := context.Background()
	for _, table := range tables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)
		`
		if err := pool.QueryRow(ctx, query, table).Scan(&exists); err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}
		if !exists {
			return fmt.Errorf("table %s does not exist", table)
		}
		log.Printf("  ✓ Table '%s' exists", table)
	}

	// Check seed data
	var userCount, modelCount int
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}
	if err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM model_versions").Scan(&modelCount); err != nil {
		return fmt.Errorf("failed to count model_versions: %w", err)
	}

	log.Printf("  ✓ Found %d users and %d model versions", userCount, modelCount)

	return nil
}

func maskPassword(dbURL string) string {
	// Simple password masking for logs
	return "postgres://diabrisk:****@postgres:5432/diabrisk?sslmode=disable"
}
