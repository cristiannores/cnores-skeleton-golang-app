package middlewares

import (
	"context"
	"github.com/labstack/echo/v4"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
	"cnores-skeleton-golang-app/app/shared/utils/log"
)

func ContextMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headerTraceParent := c.Request().Header.Get(utils_context.TraceParent)
		log.Info("header traceparent found %s", headerTraceParent)
		ctx := utils_context.CreateTraceContext(headerTraceParent)
		ctx = context.WithValue(ctx, "consumer_name", c.Request().Header.Get(""))
		log.Info("TraceID %s", ctx.Value("TraceID").(string))
		c.Set("traceContext", ctx)

		return next(c)
	}
}
