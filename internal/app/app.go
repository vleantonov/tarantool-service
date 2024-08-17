package app

import (
	"context"
	"errors"
	"fmt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tarantool/go-tarantool/v2"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"vk_tarantool_project/internal/config"
	"vk_tarantool_project/internal/handlers"
	tarantoolRepo "vk_tarantool_project/internal/repostory/tarantool"
	"vk_tarantool_project/internal/services"
)

type App struct {
	config *config.Config
	server *echo.Echo
}

// New Create app instance with new config and echo Server
// TODO: Посмотреть валидаторы
func New() *App {

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	connDB, err := getTarantoolConn(conf)
	if err != nil {
		log.Fatal(err)
	}

	repo := tarantoolRepo.New(connDB)
	service := services.New(conf, repo)
	handler := handlers.New(service)

	// Create base echo server
	server := echo.New()

	// Main server Middlewares
	server.Use(middleware.RequestID())
	server.Use(middleware.Recover())
	server.Use(middleware.Logger())
	server.Use(middleware.CORS())
	server.Use(
		middleware.RateLimiter(
			middleware.NewRateLimiterMemoryStore(10),
		),
	)

	// Data route for data requests processing
	dataRouter := server.Group("/api/")

	// JWT Middleware for dataRouter for access checking
	dataRouter.Use(echojwt.JWT([]byte(conf.Secret)))

	// Endpoints
	server.POST("/api/login", handler.Login)
	dataRouter.POST("write", nil)
	dataRouter.POST("read", nil)

	return &App{
		config: conf,
		server: server,
	}
}

// Run run application server
func (a *App) Run() {

	// Context with application terminating processing
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Start application server in gorutine
	go func() {
		if err := a.server.Start(
			fmt.Sprintf("%s:%d", a.config.Host, a.config.Port),
		); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.server.Logger.Fatal(err)
		}
	}()

	<-ctx.Done()
	a.Stop()
}

// Stop the application server
func (a *App) Stop() {

	// Context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.server.Logger.Fatal("can't shutdown app", err)
	}

	log.Println("server stopped")
}

// getTarantoolConn Connect to the Tarantool DBMS
func getTarantoolConn(conf *config.Config) (*tarantool.Connection, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  conf.TarantoolAddress,
		User:     conf.TarantoolUsername,
		Password: conf.TarantoolPassword,
	}

	opts := tarantool.Opts{
		Timeout: conf.TarantoolRequestTimeout,
	}

	return tarantool.Connect(ctx, dialer, opts)
}
