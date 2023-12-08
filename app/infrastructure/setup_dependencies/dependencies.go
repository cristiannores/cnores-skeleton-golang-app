package setup_dependencies

import (
	"context"
	"github.com/labstack/echo/v4"
	"cnores-skeleton-golang-app/app/infrastructure/constant"
	"cnores-skeleton-golang-app/app/infrastructure/setup_dependencies/dependencies"
	"cnores-skeleton-golang-app/app/infrastructure/setup_dependencies/handlers"
	"cnores-skeleton-golang-app/app/infrastructure/setup_dependencies/repositories"
	"os"

	"cnores-skeleton-golang-app/app/infrastructure/web"
	"cnores-skeleton-golang-app/app/shared/utils/config"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
	"cnores-skeleton-golang-app/app/shared/utils/log"
)

const (
	databaseOxUrl  = "database-billing-agent.url"
	databaseOxName = "database-billing-agent.name"
)

const APP = "endurance-ox-billing-agent-pro"

type InfrastructureInterface interface {
	GraceFullShutdown(sig os.Signal, log *log.Logger)
	SetupDependencies(config config.Config, echo *echo.Echo) error
}

type Infrastructure struct {
}

func NewInfrastructure() InfrastructureInterface {
	return &Infrastructure{}
}
func (infrastructure *Infrastructure) SetupDependencies(config config.Config, echo *echo.Echo) error {
	ctx := context.Background()
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "infrastructure.SetupDependencies")

	container := dependencies.NewDependencyContainer(config, echo)
	repositories.InitializeAll(container)
	handlers.InitializeAllHandlers(container)

	log.Info("Setup finished")

	// setup api sources
	web.InitRoutes()

	log.Info("Setup finished")

	return nil
}

func (infrastructure *Infrastructure) GraceFullShutdown(sig os.Signal, log *log.Logger) {
	log.Info("[App]: Init in gracefulShutdown")

	web.Shutdown()

	log.Info("[App]: Shutdown process completed for signal: %v", sig)
}
