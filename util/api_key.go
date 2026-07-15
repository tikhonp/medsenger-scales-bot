package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type agentTokenModel struct {
	AgentToken string `json:"agent_token" validate:"required"`
}

func hasRole(agentToken string, cfg *Server, allowedRoles ...string) bool {
	if agentToken == "" {
		return false
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(agentToken, claims, func(token *jwt.Token) (any, error) {
		return []byte(cfg.MedsengerAgentKey), nil
	})
	if err != nil || !token.Valid {
		return false
	}

	roles, ok := claims["roles"].([]any)
	if !ok {
		return false
	}
	for _, role := range roles {
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return true
			}
		}
	}

	return false
}

func AgentTokenJSON(cfg *Server, roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// Workaround to read request body twice
			req := c.Request()
			bodyBytes, _ := io.ReadAll(req.Body)
			if err := req.Body.Close(); err != nil {
				return err
			}
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			c.SetRequest(req)

			data := new(agentTokenModel)
			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON.")
			}
			if err := c.Validate(data); err != nil {
				return err
			}
			if !hasRole(data.AgentToken, cfg, roles...) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid agent token.")
			}
			return next(c)
		}
	}
}

func AgentTokenGetParam(cfg *Server, roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			agentToken := c.QueryParam("agent_token")
			if !hasRole(agentToken, cfg, roles...) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid agent token.")
			}
			return next(c)
		}
	}
}

func ScenarioAccess(cfg *Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			agentToken := c.Request().Header.Get("X-Agent-Token")
			if agentToken == "" {
				agentToken = c.QueryParam("agent_token")
			}
			if agentToken == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing agent token.")
			}

			if !hasRole(agentToken, cfg, "system") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid agent token.")
			}
			return next(c)
		}
	}
}
