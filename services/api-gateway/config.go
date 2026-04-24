package main

import "os"

type serviceConfig struct {
	authServiceURL string
	mlServiceURL   string
	dataServiceURL string
}

func newServiceConfig() serviceConfig {
	return serviceConfig{
		authServiceURL: getEnv("AUTH_SERVICE_URL", "http://auth-svc:8081"),
		mlServiceURL:   getEnv("ML_SERVICE_URL", "http://ml-api:8000"),
		dataServiceURL: getEnv("DATA_SERVICE_URL", "http://data-svc:8082"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
