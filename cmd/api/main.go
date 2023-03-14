package main

import (
	"github.com/alexperezortuno/go-webshell/cmd/api/bootstrap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	defer func() {
		if err := bootstrap.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}

	<-interrupt
	log.Println("Shutting down...")
}
