package handler

import (
	"database/sql"
	"fmt"
	"net/http"

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

func (h InitHandler) sendInitMessage(c db.Contract, ctx echo.Context) {
	link := fmt.Sprintf("https://scales.ai.medsenger.ru/app?agent_token=%s&contract_id=%d&type=connect", c.AgentToken.String, c.Id)
	_, err := h.MaigoClient.SendMessage(
		c.Id,
		`Если у вас есть весы Xiaomi Mi body composition scale, данные с них могут
        автоматически поступать врачу. Для этого Вам нужно скачать приложение 
        <strong>Medsenger SCALES</strong>, а затем нажать на кнопку "Подключить устройство" ниже.`,
		maigo.WithAction("Подключить устройство", link, maigo.AppUrl),
		maigo.OnlyPatient(),
	)
	if err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
}

func (h InitHandler) fetchContractDataOnInit(c db.Contract, ctx echo.Context) {
	ci, err := h.MaigoClient.GetContractInfo(c.Id)
	if err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	c.PatientName = sql.NullString{String: ci.PatientName, Valid: true}
	c.PatientEmail = sql.NullString{String: ci.PatientEmail, Valid: true}
	if err := c.Save(); err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	h.sendInitMessage(c, ctx)
	ctx.Logger().Info("Successfully fetched contract data")
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

