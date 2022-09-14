package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func UseLogMiddleware(e *echo.Echo) {
	loggerConfig := middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${uri} ${latency_human}` + "\n",
	}
	e.Use(middleware.LoggerWithConfig(loggerConfig))
}
