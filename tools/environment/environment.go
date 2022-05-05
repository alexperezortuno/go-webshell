package environment

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"time"
)

type ServerValues struct {
	Host            string
	Port            int
	ShutdownTimeout time.Duration
	Context         string
}

func env() {
	env := os.Getenv("APP_ENV")

	if env == "" || env == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}

func Server() ServerValues {
	env()
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	host := os.Getenv("APP_HOST")
	context := os.Getenv("APP_CONTEXT")

	if err != nil {
		log.Printf("error parsing port")
		port = 8080
	}

	if host == "" {
		host = "0.0.0.0"
	}

	if context == "" {
		context = "api"
	}

	return ServerValues{
		Host:            host,
		Port:            port,
		Context:         context,
		ShutdownTimeout: 10 * time.Second,
	}
}
