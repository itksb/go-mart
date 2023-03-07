package main

import (
	"context"
	"github.com/itksb/go-mart/internal/app"
	"github.com/itksb/go-mart/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg.UseOsEnv()
	cfg.UseFlags()

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	/**
	NotifyContext returns a copy of the parent context
	that is marked done (its Done channel is closed)
	 - when one of the listed signals arrives,
	 - when the returned stop function is called,
	 - or when the parent context's Done channel is closed,
	whichever happens first.

	The stop function releases resources associated with it,
	so code should call stop as soon as the operations running in this Context complete
	and signals no longer need to be diverted to the context.

	syscall.SIGTERM is the usual signal for termination
	and the default one (it can be modified) for docker containers,
	which is also used by kubernetes.
	*/
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println(application.Run(ctx))
}
