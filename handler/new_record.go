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
	BodyFatPercentage *float64       `json:"body_fat_percentage,omitempty"`
	BoneMass          *float64       `json:"bone_mass,omitempty"`
	MuscleMass        *float64       `json:"muscle_mass,omitempty"`
	WaterPercentage   *float64       `json:"water_percentage,omitempty"`
	VisceralFat       *float64       `json:"visceral_fat,omitempty"`
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
		records := []maigo.Record{
			maigo.NewRecord("weight", fmt.Sprint(m.Weight), m.Time.Time),
		}
		if m.BodyFatPercentage != nil {
			records = append(records, maigo.NewRecord("body_fat_percentage", fmt.Sprint(*m.BodyFatPercentage), m.Time.Time))
		}
		if m.BoneMass != nil {
			records = append(records, maigo.NewRecord("bone_mass", fmt.Sprint(*m.BoneMass), m.Time.Time))
		}
		if m.MuscleMass != nil {
			records = append(records, maigo.NewRecord("muscle_mass", fmt.Sprint(*m.MuscleMass), m.Time.Time))
		}
		if m.WaterPercentage != nil {
			records = append(records, maigo.NewRecord("water_percentage", fmt.Sprint(*m.WaterPercentage), m.Time.Time))
		}
		if m.VisceralFat != nil {
			records = append(records, maigo.NewRecord("visceral_fat", fmt.Sprint(*m.VisceralFat), m.Time.Time))
		}
		_, err := h.MaigoClient.AddRecords(contract.Id, records)
		if err != nil {
			sentry.CaptureException(err)
			c.Logger().Error(err)
		}
	}()

	return c.NoContent(http.StatusCreated)
}

