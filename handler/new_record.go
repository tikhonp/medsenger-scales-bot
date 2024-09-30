package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TikhonP/maigo"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-scales-bot/db"
)

type newRecordModel struct {
	AgentToken string    `json:"agent_token" validate:"required"`
	Weight     int       `json:"weight" validate:"required"`
	Time       time.Time `json:"time" validate:"required"`
}

type NewRecord struct {
	MaigoClient *maigo.Client
}

func (h NewRecord) Handle(c echo.Context) error {
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
		_, err := h.MaigoClient.AddRecord(contract.Id, "weight", fmt.Sprint(m.Weight), m.Time, nil)
		if err != nil {
			log.Println("Failed to send record to maigo:", err)
		}
		// TODO: Send record to maigo
		// so i need to implement it in maigo first :(
	}()

	return c.NoContent(http.StatusCreated)
}

