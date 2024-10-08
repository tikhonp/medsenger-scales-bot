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
	AgentToken        string         `json:"agent_token" validate:"required"`
	Time              util.Timestamp `json:"timestamp" validate:"required"`
	Weight            float64        `json:"weight" validate:"required"`
	BodyFatPercentage float64        `json:"body_fat_percentage" validate:"required"`
	BoneMass          float64        `json:"bone_mass" validate:"required"`
	MuscleMass        float64        `json:"muscle_mass" validate:"required"`
	WaterPercentage   float64        `json:"water_percentage" validate:"required"`
	VisceralFat       int            `json:"visceral_fat" validate:"required"`
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
		_, err := h.MaigoClient.AddRecords(
			contract.Id,
			[]maigo.Record{
				maigo.NewRecord("weight", fmt.Sprint(m.Weight), m.Time.Time),
				maigo.NewRecord("body_fat_percentage", fmt.Sprint(m.BodyFatPercentage), m.Time.Time),
				maigo.NewRecord("bone_mass", fmt.Sprint(m.BoneMass), m.Time.Time),
				maigo.NewRecord("muscle_mass", fmt.Sprint(m.MuscleMass), m.Time.Time),
				maigo.NewRecord("water_percentage", fmt.Sprint(m.WaterPercentage), m.Time.Time),
				maigo.NewRecord("visceral_fat", fmt.Sprint(m.VisceralFat), m.Time.Time),
			},
		)
		if err != nil {
			sentry.CaptureException(err)
			c.Logger().Error(err)
		}
	}()

	return c.NoContent(http.StatusCreated)
}

