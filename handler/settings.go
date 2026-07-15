package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type SettingsHandler struct{}

func (h SettingsHandler) Handle(c *echo.Context) error {
	return c.String(http.StatusOK, "Настройки недоступны")
}

