package db

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Contract represents Medsenger contract.
// Create on agent /init and persist during agent lifecycle.
type Contract struct {
	Id            int             `db:"id"`
	IsActive      bool            `db:"is_active"`
	AgentToken    sql.NullString  `db:"agent_token"`
	PatientName   sql.NullString  `db:"patient_name"`
	PatientEmail  sql.NullString  `db:"patient_email"`
	Locale        sql.NullString  `db:"locale"`
	PatientSex    sql.NullString  `db:"patient_sex"`
	PatientAge    sql.NullInt64   `db:"patient_age"`
	PatientHeight sql.NullFloat64 `db:"patient_height"`
}

// Save on Contract saves structure to database.
func (c *Contract) Save() error {
	const query = `
		INSERT INTO contracts (id, is_active, agent_token, patient_name, patient_email, locale, patient_sex, patient_age, patient_height)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (id)
		DO UPDATE SET is_active = EXCLUDED.is_active, agent_token = EXCLUDED.agent_token, patient_name = EXCLUDED.patient_name, patient_email = EXCLUDED.patient_email, locale = EXCLUDED.locale, patient_sex = EXCLUDED.patient_sex, patient_age = EXCLUDED.patient_age, patient_height = EXCLUDED.patient_height
	`
	_, err := db.Exec(query, c.Id, c.IsActive, c.AgentToken, c.PatientName, c.PatientEmail, c.Locale, c.PatientSex, c.PatientAge, c.PatientHeight)
	return err
}

// GetActiveContractIds returns all active contracts ids.
// Use it for medsenger status endpoint.
func GetActiveContractIds() ([]int, error) {
	var contractIds = make([]int, 0)
	err := db.Select(&contractIds, `SELECT id FROM contracts WHERE is_active = true`)
	return contractIds, err
}

// MarkInactiveContractWithId sets contract with id to inactive.
// Use it for medsenger remove endpoint.
// Equivalent to DELETE FROM contracts WHERE id = ?.
func MarkInactiveContractWithId(id int) error {
	_, err := db.Exec(`UPDATE contracts SET is_active = false WHERE id = $1`, id)
	return err
}

// GetContractById returns contract with specified id.
func GetContractById(id int) (*Contract, error) {
	contract := new(Contract)
	err := db.Get(contract, `SELECT * FROM contracts WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Contract not found")
	}
	return contract, err
}

// GetContractByAgentToken returns contract with specified agent token.
func GetContractByAgentToken(agentToken string) (*Contract, error) {
	contract := new(Contract)
	err := db.Get(contract, `SELECT * FROM contracts WHERE agent_token = $1`, agentToken)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Contract not found")
	}
	return contract, err
}

