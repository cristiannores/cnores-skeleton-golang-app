package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"cnores-skeleton-golang-app/app/infrastructure/web/routes"
	"cnores-skeleton-golang-app/app/shared/utils/config"
	"cnores-skeleton-golang-app/app/shared/utils/log"
)

var echoServer *echo.Echo

func NewWebServer() *echo.Echo {
	echoServer = echo.New()
	echoServer.HideBanner = true
	return echoServer
}

func InitRoutes() {

	echoServer.GET("/swagger/*", echoSwagger.WrapHandler)
	routes.NewHealthHandler(echoServer)
	routes.NewMetricsHandler(echoServer)
	routes.NewReadinessHandler(echoServer)
}

func Start(config config.Config) {
	log.Info("Config CurrentStage: %s", config.CurrentStage)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Port),
		ReadTimeout:  3 * time.Minute,
		WriteTimeout: 3 * time.Minute,
	}
	log.Info("App listen in port: %s", config.Port)
	echoServer.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))
	echoServer.Logger.Fatal(echoServer.StartServer(server))
}

func Shutdown() {
	log.Info("Shutting down web server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := echoServer.Shutdown(ctx); err != nil {
		log.Fatal("Error shutting down web server")
	}
}
