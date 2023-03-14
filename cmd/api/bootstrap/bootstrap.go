package bootstrap

import (
	"context"
	"github.com/alexperezortuno/go-webshell/internal/platform/server"
	"github.com/alexperezortuno/go-webshell/tools/environment"
	"log"
)

var params = environment.Server()

func Run() error {
	ctx, srv := server.New(context.Background(), params)
	return srv.Run(ctx, params)
}

func Close() error {
	log.Print("Closing server...")
	return server.Close()
}
