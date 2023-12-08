package main

import (
	"cnores-skeleton-golang-app/app/infrastructure/setup_dependencies"
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"cnores-skeleton-golang-app/app/infrastructure/constant"
	"cnores-skeleton-golang-app/app/infrastructure/web"
	"cnores-skeleton-golang-app/app/shared/utils/config"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
	_ "cnores-skeleton-golang-app/docs"
)

// @title User medical connect
// @version 1.0
// @description Microservicio de usuarios
// @host localhost:8701
// @BasePath /
// @schemes http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	ctx := context.Background()
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "Main")
	log.Info("CNORES SKELETON GOLANG APP ")
	configGlobal, err := config.GetConfig()
	if err != nil {
		log.Fatal("Error init config %s", err.Error())
	}
	//Signs Catcher
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	echoServer := web.NewWebServer()

	infra := setup_dependencies.NewInfrastructure()
	err = infra.SetupDependencies(configGlobal, echoServer)

	if err != nil {
		log.Fatal("Error setup dependencies %s", err.Error())
	}
	go web.Start(configGlobal)

	//Graceful Shutdown process
	sig := <-quit
	infra.GraceFullShutdown(sig, log)
}
