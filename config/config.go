package config

import (
	"log"
	"os"
	"strconv"
)

func GetEnv() string {
	return getEnvironmentValue("ENV")
}

func GetSqliteDB() string {
	return getEnvironmentValue("SQLITE_DB")
}

func GetApplicationGrpcPort() int {
	portStr := getEnvironmentValue("APPLICATION_GRPC_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}
	return port
}

func GetApplicationHttpPort() int {
	portStr := getEnvironmentValue("APPLICATION_HTTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}
	return port
}

func GetJwtSecret() string {
	return getEnvironmentValue("JWT_SECRET")
}

func GetJwtTokenDurationMinute() int {
	tokenDurationStr := getEnvironmentValue("JWT_TOKEN_DURATION")
	tokenDuration, err := strconv.Atoi(tokenDurationStr)
	if err != nil {
		log.Fatalf("tokenDuration: %s is invalid", tokenDurationStr)
	}
	return tokenDuration
}

func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}
	return os.Getenv(key)
}
