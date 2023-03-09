package bootstrap

import (
	"context"
	"github.com/alexperezortuno/go-webshell/internal/platform/server"
	"github.com/alexperezortuno/go-webshell/tools/environment"
)

var params = environment.Server()

func Run() error {
	ctx, srv := server.New(context.Background(), params)
	return srv.Run(ctx, params)
}
