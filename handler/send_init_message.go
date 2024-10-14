package handler

import (
	"fmt"

	"github.com/TikhonP/maigo"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-scales-bot/db"
)

func sendInitMessage(c db.Contract, maigoClient *maigo.Client, ctx echo.Context) {
	link := fmt.Sprintf(
		"https://scales.ai.medsenger.ru/app?agent_token=%s&contract_id=%d&type=connect&user_sex=%s&user_age=%d&user_height=%v",
		c.AgentToken.String, c.Id, c.PatientSex.String, c.PatientAge.Int64, c.PatientHeight.Float64,
	)
	_, err := maigoClient.SendMessage(
		c.Id,
		"Если у вас есть весы Xiaomi Mi body composition scale, данные с них могут "+
			"автоматически поступать врачу. Для этого Вам нужно скачать приложение "+
			"<strong>Medsenger Scales</strong>, а затем нажать на кнопку \"Подключить устройство\" ниже.",
		maigo.WithAction("Подключить устройство", link, maigo.AppUrl),
		maigo.OnlyPatient(),
	)
	if err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
}

