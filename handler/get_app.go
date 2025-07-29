// Package handler provides HTTP handlers for Medsenger Scales Bot. 
package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-scales-bot/util"
	"github.com/tikhonp/medsenger-scales-bot/view"
)

type GetAppHandler struct{}

func (h GetAppHandler) Handle(c echo.Context) error {
    return util.TemplRender(c, view.AppPage())
}

