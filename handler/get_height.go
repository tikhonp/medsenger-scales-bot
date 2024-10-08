package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/TikhonP/maigo"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-scales-bot/db"
	"github.com/tikhonp/medsenger-scales-bot/util"
	"github.com/tikhonp/medsenger-scales-bot/view"
)

type GetHeightHandler struct {
	MaigoClient *maigo.Client
}

func (h GetHeightHandler) Get(c echo.Context) error {
	return util.TemplRender(c, view.GetHeight(false))
}

type formData struct {
	Height int `form:"height" validate:"required"`
}

func (h GetHeightHandler) Post(c echo.Context) error {
	var fd formData
	if err := c.Bind(&fd); err != nil {
		return err
	}
	if err := c.Validate(fd); err != nil {
		return err
	}
	contractId := util.QueryParamInt(c, "contract_id")
	if contractId == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "contract_id is required")
	}
	contract, err := db.GetContractById(*contractId)
	if err != nil {
		return err
	}
	contract.PatientHeight = sql.NullFloat64{Float64: float64(fd.Height), Valid: true}
	if err := contract.Save(); err != nil {
		return err
	}
	go func() {
        sendInitMessage(*contract, h.MaigoClient, c)
		_, err := h.MaigoClient.AddRecords(*contractId, []maigo.Record{maigo.NewRecord("height", strconv.Itoa(fd.Height), time.Now())})
		if err != nil {
			sentry.CaptureException(err)
			c.Logger().Error(err)
		}
	}()
	return util.TemplRender(c, view.GetHeight(true))
}

