// Package medsengerscalesbot provides the main server for the Medsenger Scales Bot.
package medsengerscalesbot

import (
	"fmt"
	"os"
	"time"

	"github.com/tikhonp/maigo"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tikhonp/medsenger-scales-bot/handler"
	"github.com/tikhonp/medsenger-scales-bot/util"
)

type Server struct {
	cfg       *util.Server
	root      handler.RootHandler
	init      handler.InitHandler
	status    handler.StatusHandler
	remove    handler.RemoveHandler
	settings  handler.SettingsHandler
	scenario  handler.ScenarioCapabilitiesHandler
	newRecord handler.NewRecordHandler
	getApp    handler.GetAppHandler
	getHeight handler.GetHeightHandler
}

func NewServer(cfg *util.Server) *Server {
	maigoClient := maigo.Init(cfg.MedsengerAgentKey)
	return &Server{
		cfg:       cfg,
		init:      handler.InitHandler{MaigoClient: maigoClient},
		newRecord: handler.NewRecordHandler{MaigoClient: maigoClient},
		getHeight: handler.GetHeightHandler{MaigoClient: maigoClient},
	}
}

func (s *Server) Listen() {
	app := echo.New()
	app.Validator = util.NewDefaultValidator()

	if !s.cfg.Debug {
		app.Use(sentryecho.New(sentryecho.Options{
			Repanic:         true,
			WaitForDelivery: false,
			Timeout:         5 * time.Second,
		}))
	}
	app.Use(middleware.RequestLogger())
	app.Use(middleware.Recover())

	app.File("/.well-known/apple-app-site-association", "public/apple-app-site-association.json")
	app.File("/.well-known/assetlinks.json", "public/assetlinks.json")
	app.Static("/static", "public/static")
	app.GET("/", s.root.Handle)
	app.POST("/init", s.init.Handle, util.AgentTokenJSON(s.cfg, "system"))
	app.POST("/status", s.status.Handle, util.AgentTokenJSON(s.cfg, "system"))
	app.POST("/remove", s.remove.Handle, util.AgentTokenJSON(s.cfg, "system"))
	app.GET("/scenario-capabilities/v1", s.scenario.Capabilities, util.ScenarioAccess(s.cfg))
	app.GET("/scenario-capabilities/v1/objects/:object_type", s.scenario.Objects, util.ScenarioAccess(s.cfg))
	app.GET("/scenario-capabilities/v1/objects/:object_type/:object_id", s.scenario.Object, util.ScenarioAccess(s.cfg))
	app.GET("/settings", s.settings.Handle, util.AgentTokenGetParam(s.cfg, "doctor", "patient", "system"))
	app.POST("/new_record", s.newRecord.Handle)
	app.GET("/app", s.getApp.Handle)

	app.GET("/get_height", s.getHeight.Get, util.AgentTokenGetParam(s.cfg, "doctor", "patient", "system"))
	app.POST("/get_height", s.getHeight.Post, util.AgentTokenGetParam(s.cfg, "doctor", "patient", "system"))

	addr := fmt.Sprintf(":%d", s.cfg.Port)
	if err := app.Start(addr); err != nil {
		app.Logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
