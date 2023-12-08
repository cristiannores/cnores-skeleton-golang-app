package api_routes

import (
	"github.com/labstack/echo/v4"
	"cnores-skeleton-golang-app/app/infrastructure/web/middlewares"
	"cnores-skeleton-golang-app/app/interfaces/input/handlers"
	"cnores-skeleton-golang-app/app/interfaces/input/handlers/service/save_service_handler"
	"cnores-skeleton-golang-app/app/interfaces/input/handlers/service/search_services_handler"
	"cnores-skeleton-golang-app/app/interfaces/input/handlers/service/update_service_handler"
)

func NewServiceHandler(e *echo.Echo,
	getPatientHandlerDetail handlers.IGetServiceByIdRestHandler,
	searchServicesHandler search_services_handler.ISearchServicesHandler,
	saveServiceHandler save_service_handler.ISaveServiceHandler,
	updateServiceHandler update_service_handler.IUpdateServiceHandler,
) {
	e.GET("/api/services/:id", getPatientHandlerDetail.Handle, middlewares.ContextMiddleWare)
	e.POST("/api/services/search", searchServicesHandler.Handle, middlewares.ContextMiddleWare)
	e.POST("/api/services", saveServiceHandler.Handle, middlewares.ContextMiddleWare)
	e.PATCH("/api/services/:id", updateServiceHandler.Handle, middlewares.ContextMiddleWare)
}
