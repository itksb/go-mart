package app

import (
	"errors"
	"github.com/itksb/go-mart/internal/config"
	"github.com/itksb/go-mart/internal/handler"
	"github.com/itksb/go-mart/internal/router"
	"github.com/itksb/go-mart/internal/service/auth"
	"github.com/itksb/go-mart/internal/service/auth/token"
	"github.com/itksb/go-mart/internal/storage/dbpgsql"
	"github.com/itksb/go-mart/pkg/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type App struct {
	HTTPServer *http.Server
	logger     logger.Interface
	auth       *auth.Service

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

	db, err := dbpgsql.NewPostgres(cfg.DatabaseURI, sugar)
	if err != nil {
		sugar.Errorf("unable connect to postgres: %s", err.Error())
		return nil, err
	}

	identityProvider, err := dbpgsql.NewIdentityPostgres(db)
	if err != nil {
		sugar.Errorf("unable create identity provider: %s", err.Error())
		return nil, err
	}

	hashAlgo, err := auth.NewHashAlgoBcrypt()
	if err != nil {
		sugar.Errorf("unable initialize hash algo module: %s", err.Error())
		return nil, err
	}

	authService, err := auth.NewAuthService(auth.Opts{
		IdentityProvider: identityProvider,
		Logger:           sugar,
		HashAlgo:         hashAlgo,
		TokenCreate:      token.CreateToken,
		TokenParse:       token.ParseWithClaims,
		SecretReader: token.SecretFunc(func() (string, error) {
			return cfg.AppSecret, nil
		}),
		NowTime: time.Now,
	})

	if err != nil {
		sugar.Errorf("unable to initialize auth service: %s", err.Error())
		return nil, err
	}

	return &App{
		HTTPServer: srv,
		logger:     sugar,
		auth:       authService,
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
