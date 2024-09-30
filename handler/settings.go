package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SettingsHandler struct{}

func (h SettingsHandler) Handle(c echo.Context) error {
	return c.String(http.StatusOK, "Настройки недоступны")
}

