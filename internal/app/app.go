package app

import (
	"context"
	"errors"
	"github.com/itksb/go-mart/internal/config"
	"github.com/itksb/go-mart/internal/handler"
	"github.com/itksb/go-mart/internal/router"
	"github.com/itksb/go-mart/internal/service/auth"
	"github.com/itksb/go-mart/internal/service/auth/token"
	"github.com/itksb/go-mart/internal/service/balance"
	"github.com/itksb/go-mart/internal/service/order"
	"github.com/itksb/go-mart/internal/service/withdraw"
	"github.com/itksb/go-mart/internal/storage/pg"
	"github.com/itksb/go-mart/internal/worker"
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
	orderService             *order.Service
	withdrawService          *withdraw.Service
	balanceService           *balance.Service
	onCloseCallbacks         []OnCloseCallback
	db                       *sqlx.DB
	worker                   *worker.AccrualWorker
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
func (a *App) Run(sigCtx context.Context) error {

	// run migrations (from embedded fs)
	err := migrate.Migrate(a.dsn, migrate.Migrations)
	if err != nil {
		a.logger.Errorf("migration error: %s", err.Error())
		return err
	}

	err = a.connectDB(sigCtx)
	if err != nil {
		return err
	}

	err = a.setupAuthService()
	if err != nil {
		return err
	}
	a.logger.Infof("auth service created")

	err = a.setupDomainServices()
	if err != nil {
		return err
	}
	a.logger.Infof("domain services created")

	workerStorage, err := worker.NewWorkerStorage(a.db)
	if err != nil {
		return err
	}
	accWorker, err := worker.NewAccrualWorker(a.cfg.AccrualSystemAddress, workerStorage, a.logger)
	if err != nil {
		return err
	}
	a.worker = accWorker

	err = a.setupServer()
	if err != nil {
		return err
	}
	a.logger.Infof("http server created")

	group, groupCtx := errgroup.WithContext(sigCtx)
	// run the server
	group.Go(func() error {
		a.logger.Infof("server is starting: %s", a.httpServer.Addr)
		/**
		When Shutdown is called, Serve, ListenAndServe, and
		ListenAndServeTLS immediately return ErrServerClosed. Make sure the
		program doesn't exit and waits instead for Shutdown to return.
		*/
		err := a.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Errorf("failed to start server at %s. Error: %s", a.httpServer.Addr, err.Error())
		}
		return err
	})

	group.Go(func() error {
		return a.worker.Run(groupCtx)
	})

	// graceful shutdown the web server
	group.Go(func() error {
		<-groupCtx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(a.gracefulShutdownInterval)*time.Second)
		defer cancel()
		a.logger.Infof("shutting down the server...")
		err := a.httpServer.Shutdown(timeoutCtx)
		if err != nil {
			a.logger.Errorf("failed to shutdown the server %s. Error: %s", a.httpServer.Addr, err.Error())
		}
		return err
	})

	defer a.onClose()
	return group.Wait()
}

// Close -
func (a *App) Close() error {
	a.onClose()
	return nil
}

func (a *App) onClose() {
	if a.closed {
		return
	}
	for _, clbk := range a.onCloseCallbacks {
		clbk(a)
	}
	a.closed = true
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

func (a *App) setupDomainServices() error {
	orderRep, err := pg.NewOrderRepositoryPg(a.db)
	if err != nil {
		return err
	}
	a.orderService, err = order.NewOrderService(orderRep)
	if err != nil {
		return err
	}

	balanceRep, err := pg.NewBalanceRepositoryPg(a.db)
	if err != nil {
		return err
	}
	a.balanceService, err = balance.NewBalanceService(balanceRep)
	if err != nil {
		return err
	}

	withdrawRep, err := pg.NewWithdrawRepositoryPg(a.db)
	if err != nil {
		return err
	}

	a.withdrawService, err = withdraw.NewWithdrawService(withdrawRep, balanceRep)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) setupServer() error {
	if a.cfg.AppHost == "" {
		return config.ErrConfigAppHostIsEmpty
	}

	h := handler.NewHandler(
		a.logger,
		*a.cfg,
		a.auth,
		a.orderService,
		a.withdrawService,
		a.balanceService,
	)
	routeHandler, err := router.NewRouter(h, a.logger, a.auth)
	if err != nil {
		a.logger.Errorf("Router creating error: %s", err.Error())
		return err
	}

	srv := &http.Server{
		Addr:         a.cfg.GetFullAddr(),
		Handler:      routeHandler,
		WriteTimeout: 15 * time.Second,
	}

	a.httpServer = srv

	return nil
}

func (a *App) connectDB(ctx context.Context) error {
	if a.dsn == "" {
		return config.ErrConfigDatabaseURIEmpty
	}

	// this Pings the database trying to connect
	db, err := sqlx.ConnectContext(ctx, "postgres", a.dsn)
	if err != nil {
		a.logger.Errorf("database connection error %s", err.Error())
		return err
	}
	a.AddOnCloseCallback(func(a *App) { a.db.Close() })
	a.db = db

	return nil
}

func (a *App) setupAuthService() error {
	identityProvider, err := pg.NewIdentityProviderPg(a.db)
	if err != nil {
		a.logger.Errorf("unable create identity provider: %s", err.Error())
		return err
	}

	hashAlgo, err := auth.NewHashAlgoBcrypt()
	if err != nil {
		a.logger.Errorf("unable initialize hash algo module: %s", err.Error())
		return err
	}

	authService, err := auth.NewAuthService(auth.Opts{
		IdentityProvider: identityProvider,
		Logger:           a.logger,
		HashAlgo:         hashAlgo,
		TokenCreate: func(martClaims *token.MartClaims, secretReader token.Secret) (newToken string, err error) {
			return token.CreateToken(martClaims, secretReader, time.Now)
		},
		TokenParse: token.ParseWithClaims,
		SecretReader: token.SecretFunc(func() (string, error) {
			return a.appSecret, nil
		}),
		NowTime: time.Now,
	})

	if err != nil {
		a.logger.Errorf("unable to initialize auth service: %s", err.Error())
		return err
	}

	a.auth = authService

	return nil
}
