package app

import (
	"context"
	"errors"
	"github.com/itksb/go-mart/internal/config"
	"github.com/itksb/go-mart/internal/handler"
	"github.com/itksb/go-mart/internal/router"
	"github.com/itksb/go-mart/internal/service/auth"
	"github.com/itksb/go-mart/internal/service/auth/token"
	"github.com/itksb/go-mart/internal/storage/pg"
	"github.com/itksb/go-mart/migrate"
	"github.com/itksb/go-mart/pkg/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"time"

	//Under the hood, the driver registers itself as being available to the database/sql package,
	//but in general nothing else happens with the exception that the init function is run.
	_ "github.com/lib/pq"
)

type OnCloseCallback func(a *App)

type App struct {
	httpServer               *http.Server
	logger                   logger.Interface
	auth                     *auth.Service
	onCloseCallbacks         []OnCloseCallback
	db                       *sqlx.DB
	gracefulShutdownInterval time.Duration //seconds
	closed                   bool
	dsn                      string
	appSecret                string
	cfg                      *config.Config

	io.Closer
}

func (a *App) AddOnCloseCallback(clbck OnCloseCallback) {
	a.onCloseCallbacks = append(a.onCloseCallbacks, clbck)
}

func NewApp(cfg config.Config) (*App, error) {
	app := &App{
		onCloseCallbacks: make([]OnCloseCallback, 0),
	}
	defer app.onClose()

	err := app.setupLogger(&cfg)
	if err != nil {
		return nil, err
	}
	app.logger.Infof("logger created")

	app.gracefulShutdownInterval = cfg.GracefulShutdownInterval
	app.closed = false

	if cfg.DatabaseURI == "" {
		return nil, config.ErrConfigDatabaseURIEmpty
	}
	app.dsn = cfg.DatabaseURI
	app.appSecret = cfg.AppSecret

	app.cfg = &cfg

	return app, nil
}

// Run - run the application instance
func (app *App) Run(sigCtx context.Context) error {

	// run migrations (from embedded fs)
	err := migrate.Migrate(app.dsn, migrate.Migrations)
	if err != nil {
		app.logger.Errorf("migration error: %s", err.Error())
		return err
	}

	err = app.connectDb(sigCtx)
	if err != nil {
		return err
	}

	err = app.setupAuthService()
	if err != nil {
		return err
	}
	app.logger.Infof("auth service created")

	err = app.setupServer()
	if err != nil {
		return err
	}
	app.logger.Infof("http server created")

	group, groupCtx := errgroup.WithContext(sigCtx)
	// run the server
	group.Go(func() error {
		app.logger.Infof("server is starting: %s", app.httpServer.Addr)
		/**
		When Shutdown is called, Serve, ListenAndServe, and
		ListenAndServeTLS immediately return ErrServerClosed. Make sure the
		program doesn't exit and waits instead for Shutdown to return.
		*/
		err := app.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Errorf("failed to start server at %s. Error: %s", app.httpServer.Addr, err.Error())
		}
		return err
	})

	// graceful shutdown the web server
	group.Go(func() error {
		<-groupCtx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(app.gracefulShutdownInterval)*time.Second)
		defer cancel()
		app.logger.Infof("shutting down the server...")
		err := app.httpServer.Shutdown(timeoutCtx)
		if err != nil {
			app.logger.Errorf("failed to shutdown the server %s. Error: %s", app.httpServer.Addr, err.Error())
		}
		return err
	})

	defer app.onClose()
	return group.Wait()
}

// Close -
func (app *App) Close() error {
	app.onClose()
	return nil
}

func (app *App) onClose() {
	if app.closed {
		return
	}
	for _, clbk := range app.onCloseCallbacks {
		clbk(app)
	}
	app.closed = true
}

func (a *App) setupLogger(cfg *config.Config) error {
	var appLogger *zap.Logger
	var err error

	switch cfg.Env {
	case config.Debug:
		appLogger, err = zap.NewDevelopment()
	case config.Prod:
		appLogger, err = zap.NewProduction()
	default:
		return config.ErrConfigWrongEnvValue
	}

	a.logger = appLogger.Sugar()

	a.AddOnCloseCallback(func(a *App) {
		a.logger.(*zap.SugaredLogger).Sync()
	})

	return err

}

func (app *App) setupServer() error {
	if app.cfg.AppHost == "" {
		return config.ErrConfigAppHostIsEmpty
	}

	h := handler.NewHandler(app.logger, *app.cfg, app.auth)
	routeHandler, err := router.NewRouter(h, app.logger, app.auth)
	if err != nil {
		app.logger.Errorf("Router creating error: %s", err.Error())
		return err
	}

	srv := &http.Server{
		Addr:         app.cfg.GetFullAddr(),
		Handler:      routeHandler,
		WriteTimeout: 15 * time.Second,
	}

	app.httpServer = srv

	return nil
}

func (app *App) connectDb(ctx context.Context) error {
	if app.dsn == "" {
		return config.ErrConfigDatabaseURIEmpty
	}

	// this Pings the database trying to connect
	db, err := sqlx.ConnectContext(ctx, "postgres", app.dsn)
	if err != nil {
		app.logger.Errorf("database connection error %s", err.Error())
		return err
	}
	app.AddOnCloseCallback(func(a *App) { a.db.Close() })
	app.db = db

	return nil
}

func (app *App) setupAuthService() error {
	identityProvider, err := pg.NewIdentityProviderPg(app.db)
	if err != nil {
		app.logger.Errorf("unable create identity provider: %s", err.Error())
		return err
	}

	hashAlgo, err := auth.NewHashAlgoBcrypt()
	if err != nil {
		app.logger.Errorf("unable initialize hash algo module: %s", err.Error())
		return err
	}

	authService, err := auth.NewAuthService(auth.Opts{
		IdentityProvider: identityProvider,
		Logger:           app.logger,
		HashAlgo:         hashAlgo,
		TokenCreate: func(martClaims *token.MartClaims, secretReader token.Secret) (newToken string, err error) {
			return token.CreateToken(martClaims, secretReader, time.Now)
		},
		TokenParse: token.ParseWithClaims,
		SecretReader: token.SecretFunc(func() (string, error) {
			return app.appSecret, nil
		}),
		NowTime: time.Now,
	})

	if err != nil {
		app.logger.Errorf("unable to initialize auth service: %s", err.Error())
		return err
	}

	app.auth = authService

	return nil
}
