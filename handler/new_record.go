package handler

import (
	"fmt"
	"net/http"

	"github.com/TikhonP/maigo"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-scales-bot/db"
	"github.com/tikhonp/medsenger-scales-bot/util"
)

type newRecordModel struct {
	AgentToken string         `json:"agent_token" validate:"required"`
	Weight     int            `json:"weight" validate:"required"`
	Time       util.Timestamp `json:"timestamp" validate:"required"`
	OtherData  string         `json:"other_data" validate:"required"`
}

type NewRecordHandler struct {
	MaigoClient *maigo.Client
}

func (h NewRecordHandler) Handle(c echo.Context) error {
	m := new(newRecordModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}
	contract, err := db.GetContractByAgentToken(m.AgentToken)
	if err != nil {
		return err
	}

	go func() {
		_, err := h.MaigoClient.AddRecords(contract.Id, []maigo.Record{maigo.NewRecord("weight", fmt.Sprint(m.Weight), m.Time.Time)})
		if err != nil {
			sentry.CaptureException(err)
			c.Logger().Error(err)
		}
		_, err = h.MaigoClient.AddRecords(contract.Id, []maigo.Record{maigo.NewRecord("information", m.OtherData, m.Time.Time)})
		if err != nil {
			sentry.CaptureException(err)
			c.Logger().Error(err)
		}
	}()

	return c.NoContent(http.StatusCreated)
}

