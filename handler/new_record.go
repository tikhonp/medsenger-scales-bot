package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TikhonP/maigo"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-scales-bot/db"
	"github.com/tikhonp/medsenger-scales-bot/util"
)

type newRecordModel struct {
	AgentToken string         `json:"agent_token" validate:"required"`
	Weight     int            `json:"weight" validate:"required"`
	Time       util.Timestamp `json:"timestamp" validate:"required"`
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
			log.Println("Failed to send record to maigo:", err)
		}
		// TODO: Send record to maigo
		// so i need to implement it in maigo first :(
	}()

	return c.NoContent(http.StatusCreated)
}

