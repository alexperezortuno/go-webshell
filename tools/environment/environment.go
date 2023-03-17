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
	SafetyZone      string
	SecretSession   string
	SessionName     string
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
	value := os.Getenv(envName)
	res, err := strconv.Atoi(value)
	if err != nil {
		return valueDefault
	}
	return res
}

func Server() ServerValues {
	env()
	port := getEnvInt("APP_PORT", 8001)
	host := getEnv("APP_HOST", "0.0.0.0")
	timeZone := getEnv("APP_TIME_ZONE", "America/Santiago")
	context := getEnv("APP_CONTEXT", "api")
	safetyZone := getEnv("SAFETY_ZONE", "/safety_zone")
	secretSession := getEnv("SECRET_SESSION", "secret")
	sessionName := getEnv("SESSION_NAME", "session")

	return ServerValues{
		Host:            host,
		Port:            port,
		Context:         context,
		TimeZone:        timeZone,
		ShutdownTimeout: 10 * time.Second,
		SafetyZone:      safetyZone,
		SecretSession:   secretSession,
		SessionName:     sessionName,
	}
}
