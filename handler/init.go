package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/TikhonP/maigo"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-scales-bot/db"
)

type initModel struct {
	ContractId        int    `json:"contract_id" validate:"required"`
	ClinicId          int    `json:"clinic_id" validate:"required"`
	AgentToken        string `json:"agent_token" validate:"required"`
	PatientAgentToken string `json:"patient_agent_token" validate:"required"`
	DoctorAgentToken  string `json:"doctor_agent_token" validate:"required"`
	AgentId           int    `json:"agent_id" validate:"required"`
	AgentName         string `json:"agent_name" validate:"required"`
	Locale            string `json:"locale" validate:"required"`
}

type InitHandler struct {
	MaigoClient *maigo.Client
}

func (h InitHandler) fetchContractDataOnInit(c db.Contract, ctx echo.Context) {
	ci, err := h.MaigoClient.GetContractInfo(c.Id)
	if err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	age, err := strconv.ParseInt(ci.PatientAge, 10, 64)
	if err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	records, err := h.MaigoClient.GetRecords(c.Id, maigo.WithCategoryName("height"), maigo.Limit(1))
	if err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	var height float64
	if len(records) == 0 {
		castedHeight, ok := records[0].Value.(float64)
        if ok {
            height = castedHeight
        }
	}
	if ci.PatientSex == "" {
		c.PatientName = sql.NullString{Valid: false}
	} else {
		c.PatientName = sql.NullString{String: ci.PatientName, Valid: true}
	}
	if ci.PatientEmail == "" {
		c.PatientEmail = sql.NullString{Valid: false}
	} else {
		c.PatientEmail = sql.NullString{String: ci.PatientEmail, Valid: true}
	}
	if ci.PatientSex == "" {
		c.PatientSex = sql.NullString{Valid: false}
	} else {
		c.PatientSex = sql.NullString{String: string(ci.PatientSex), Valid: true}
	}
	if ci.PatientAge == "" {
		c.PatientAge = sql.NullInt64{Valid: false}
	} else {
		c.PatientAge = sql.NullInt64{Int64: age, Valid: true}
	}
	if height == 0 {
		c.PatientHeight = sql.NullFloat64{Valid: false}
	} else {
		c.PatientHeight = sql.NullFloat64{Float64: height, Valid: true}
	}
	if err := c.Save(); err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	if c.PatientHeight.Valid {
		sendInitMessage(c, h.MaigoClient, ctx)
	} else {
		_, err := h.MaigoClient.SendMessage(
			c.Id,
			"Пожалуйста, введите свой рост, чтобы мы могли рассчитать индекс массы Вашего тела.",
			maigo.WithAction("Заполнить", "/get_height", maigo.Action),
			maigo.OnlyPatient(),
		)
		if err != nil {
			sentry.CaptureException(err)
			ctx.Logger().Error(err)
		}
	}
}

func (h InitHandler) Handle(c echo.Context) error {
	m := new(initModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}
	contract := db.Contract{
		Id:         m.ContractId,
		IsActive:   true,
		AgentToken: sql.NullString{String: m.AgentToken, Valid: true},
		Locale:     sql.NullString{String: m.Locale, Valid: true},
	}
	if err := contract.Save(); err != nil {
		return err
	}
	go h.fetchContractDataOnInit(contract, c)
	return c.String(http.StatusCreated, "ok")
}

