package medsengerscalesbot

import (
	"fmt"

	"github.com/TikhonP/maigo"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tikhonp/medsenger-scales-bot/config"
	"github.com/tikhonp/medsenger-scales-bot/handler"
	"github.com/tikhonp/medsenger-scales-bot/util"
)

type Server struct {
	cfg       *config.Server
	root      handler.RootHandler
	init      handler.InitHandler
	status    handler.StatusHandler
	remove    handler.RemoveHandler
	settings  handler.SettingsHandler
	newRecord handler.NewRecordHandler
	getApp    handler.GetAppHandler
	getHeight handler.GetHeightHandler
}

func NewServer(cfg *config.Server) *Server {
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
	app.Debug = s.cfg.Debug
	app.HideBanner = true
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
		Output: app.Logger.Output(),
	}))
	app.Use(middleware.Recover())
	if !s.cfg.Debug {
		app.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
		app.Logger.Printf("Sentry initialized")
	}
	app.Validator = util.NewDefaultValidator()

	app.File("/.well-known/apple-app-site-association", "public/apple-app-site-association.json")
	app.File("/.well-known/assetlinks.json", "public/assetlinks.json")
	app.Static("/static", "public/static")
	app.GET("/", s.root.Handle)
	app.POST("/init", s.init.Handle, util.ApiKeyJSON(s.cfg))
	app.POST("/status", s.status.Handle, util.ApiKeyJSON(s.cfg))
	app.POST("/remove", s.remove.Handle, util.ApiKeyJSON(s.cfg))
	app.GET("/settings", s.settings.Handle, util.ApiKeyGetParam(s.cfg))
	app.POST("/new_record", s.newRecord.Handle)
	app.GET("/app", s.getApp.Handle)

	app.GET("/get_height", s.getHeight.Get, util.ApiKeyGetParam(s.cfg))
	app.POST("/get_height", s.getHeight.Post, util.ApiKeyGetParam(s.cfg))

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	app.Logger.Fatal(app.Start(addr))
}

