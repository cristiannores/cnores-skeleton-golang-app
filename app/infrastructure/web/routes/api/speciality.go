package api_routes

import (
	"github.com/labstack/echo/v4"
	"cnores-skeleton-golang-app/app/infrastructure/web/middlewares"
	"cnores-skeleton-golang-app/app/interfaces/input/handlers/speciality/get_speciality_by_id_handler"
	"cnores-skeleton-golang-app/app/interfaces/input/handlers/speciality/save_speciality_handler"
	"cnores-skeleton-golang-app/app/interfaces/input/handlers/speciality/search_specialities_handler"
	"cnores-skeleton-golang-app/app/interfaces/input/handlers/speciality/update_speciality_handler"
)

func NewSpecialityHandler(e *echo.Echo,
	getSpecialityByIDHandler get_speciality_by_id_handler.IGetSpecialityByIdHandler,
	searchSpecialityHandler search_specialities_handler.ISearchSpecialitiesHandler,
	saveSpecialityHandler save_speciality_handler.ISaveSpecialityHandler,
	updateSpecialityHandler update_speciality_handler.IUpdateSpecialityHandler,
) {
	e.GET("/api/specialities/:specialityID", getSpecialityByIDHandler.Handle, middlewares.ContextMiddleWare)
	e.POST("/api/specialities/search", searchSpecialityHandler.Handle, middlewares.ContextMiddleWare)
	e.POST("/api/specialities", saveSpecialityHandler.Handle, middlewares.ContextMiddleWare)
	e.PATCH("/api/specialities/:id", updateSpecialityHandler.Handle, middlewares.ContextMiddleWare)
}
