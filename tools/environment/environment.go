package environment

import (
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"time"
)

type ServerValues struct {
	Host            string
	Port            int
	TimeZone        string
	ShutdownTimeout time.Duration
	Context         string
}

func env() {
	env := os.Getenv("APP_ENV")

	if env == "" || env == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func getEnv(envName, valueDefault string) string {
	value := os.Getenv(envName)
	if value == "" {
		return valueDefault
	}
	return value
}

func getEnvInt(envName string, valueDefault int) int {
	value, err := strconv.Atoi(envName)
	if err != nil {
		return valueDefault
	}
	return value
}

func Server() ServerValues {
	env()
	port := getEnvInt("APP_PORT", 8080)
	host := getEnv("APP_HOST", "0.0.0.0")
	timeZone := getEnv("APP_TIME_ZONE", "America/Santiago")
	context := getEnv("APP_CONTEXT", "api")

	return ServerValues{
		Host:            host,
		Port:            port,
		Context:         context,
		TimeZone:        timeZone,
		ShutdownTimeout: 10 * time.Second,
	}
}
