package main

import "os"

type serviceConfig struct {
	dbURL string
	port  string
}

func loadConfig() serviceConfig {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://diabrisk:diabrisk123@postgres:5432/diabrisk?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	return serviceConfig{
		dbURL: dbURL,
		port:  port,
	}
}

func maskPassword(dbURL string) string {
	// Simple password masking for logs
	return "postgres://diabrisk:****@postgres:5432/diabrisk?sslmode=disable"
}
