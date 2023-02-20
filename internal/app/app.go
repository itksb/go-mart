package app

import (
	"errors"
	"github.com/itksb/go-mart/internal/config"
	"github.com/itksb/go-mart/internal/handler"
	"github.com/itksb/go-mart/internal/router"
	"github.com/itksb/go-mart/pkg/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type App struct {
	HTTPServer *http.Server
	logger     logger.Interface

	io.Closer
}

func NewApp(cfg config.Config) (*App, error) {
	if cfg.AppHost == "" {
		return nil, errors.New("Wrong configuration. AppHost is empty.")
	}
	if cfg.DatabaseURI == "" {
		return nil, errors.New("Wrong configuration. DatabaseURI is empty.")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infof("zap logger created")

	h := handler.NewHandler(sugar, cfg)

	routeHandler, err := router.NewRouter(h, sugar)
	if err != nil {
		sugar.Errorf("Router creating error: %s", err.Error())
		return nil, err
	}

	srv := &http.Server{
		Addr:         cfg.GetFullAddr(),
		Handler:      routeHandler,
		WriteTimeout: 15 * time.Second,
	}

	return &App{
		HTTPServer: srv,
		logger:     sugar,
	}, nil
}

// Run - run the application instance
func (app *App) Run() error {
	app.logger.Infof("server started: %s", app.HTTPServer.Addr)
	return app.HTTPServer.ListenAndServe()
}

// Close -
func (app *App) Close() error {
	return nil
}
