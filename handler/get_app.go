package handler

import "github.com/labstack/echo/v4"

type GetAppHandler struct{}

func (h GetAppHandler) Handle(c echo.Context) error {
	return c.File("public/get_app.html")
}

