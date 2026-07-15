package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

const scalesScenarioObjectType = "scales_devices"

type ScenarioCapabilitiesHandler struct{}

func scalesScenarioItem() map[string]any {
	return map[string]any{
		"id":          1,
		"title":       "Xiaomi Mi Scale",
		"description": "Умные весы Xiaomi для мониторинга веса и состава тела.",
	}
}

func (h ScenarioCapabilitiesHandler) Capabilities(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"version": 1,
		"agent": map[string]any{
			"code":  "scales",
			"title": "Умные весы",
		},
		"params": []any{},
		"object_types": []map[string]any{
			{
				"type":        scalesScenarioObjectType,
				"title":       "Умные весы",
				"description": "Подключение весов Xiaomi Mi Scale.",
				"icon":        "scale",
				"selection":   "many",
			},
		},
	})
}

func (h ScenarioCapabilitiesHandler) Objects(c *echo.Context) error {
	if c.Param("object_type") != scalesScenarioObjectType {
		return echo.NewHTTPError(http.StatusNotFound, "Object type not found.")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"type":  scalesScenarioObjectType,
		"items": []map[string]any{scalesScenarioItem()},
	})
}

func (h ScenarioCapabilitiesHandler) Object(c *echo.Context) error {
	if c.Param("object_type") != scalesScenarioObjectType || c.Param("object_id") != "1" {
		return echo.NewHTTPError(http.StatusNotFound, "Object not found.")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"type":   scalesScenarioObjectType,
		"object": scalesScenarioItem(),
		"params": []any{},
	})
}
