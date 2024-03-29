package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alexperezortuno/go-webshell/internal/platform/server/handler/command"
	"github.com/alexperezortuno/go-webshell/internal/platform/server/handler/health"
	"github.com/alexperezortuno/go-webshell/internal/platform/server/handler/instance"
	"github.com/alexperezortuno/go-webshell/internal/platform/server/middleware/logging"
	"github.com/alexperezortuno/go-webshell/internal/platform/server/middleware/recovery"
	"github.com/alexperezortuno/go-webshell/internal/platform/storage/data_base"
	"github.com/alexperezortuno/go-webshell/tools/environment"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	httpAddr        string
	engine          *gin.Engine
	shutdownTimeout time.Duration
}

var p = environment.Server()

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}

func New(ctx context.Context, params environment.ServerValues) (context.Context, Server) {
	srv := Server{
		engine:          gin.New(),
		httpAddr:        fmt.Sprintf("%s:%d", params.Host, params.Port),
		shutdownTimeout: params.ShutdownTimeout,
	}

	log.Println(fmt.Sprintf("Check app in %s:%d/%s/%s", params.Host, params.Port, params.Context, "health"))
	srv.registerRoutes(params.Context)
	return serverContext(ctx), srv
}

func Close() error {
	log.Print("Closing db connection...")
	err := data_base.Close()
	if err != nil {
		log.Printf("Error closing db: %s", err)
		return err
	}

	return nil
}

func (s *Server) Run(ctx context.Context, params environment.ServerValues) error {
	log.Println("Server running on", s.httpAddr)
	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	db, err := data_base.Connect()
	if db == nil {
		log.Fatalf("No connection to database: %s", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing db: %s", err)
		}
	}(db)

	_, err = data_base.Start()
	if err != nil {
		log.Fatalf("Error starting db: %s", err)
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func (s *Server) registerRoutes(context string) {
	store := memstore.NewStore([]byte(p.SecretSession))
	s.engine.Use(sessions.Sessions(p.SessionName, store))
	s.engine.Use(gzip.Gzip(gzip.DefaultCompression))
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	s.engine.Use(recovery.Middleware(), logging.Middleware())
	s.engine.Use(logging.Middleware(), gin.Logger(), recovery.Middleware())
	s.engine.GET(fmt.Sprintf("/%s/%s", context, "/health"), health.CheckHandler())
	s.engine.GET(fmt.Sprintf("/%s/%s", context, "/shell"), command.CmdHandler())
	s.engine.GET(fmt.Sprintf("/%s/%s", context, "/instance"), instance.Handler())
}
