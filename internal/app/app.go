package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"vk_tarantool_project/internal/config"
)

type App struct {
	config *config.Config
	server *echo.Echo
}

// New Create app instance with new config and echo Server
func New() *App {

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	server := echo.New()

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
