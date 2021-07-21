package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a App) HealthCheck(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "healthy")
}