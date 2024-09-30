package medsengerscalesbot

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tikhonp/medsenger-scales-bot/config"
	"github.com/tikhonp/medsenger-scales-bot/handler"
	"github.com/tikhonp/medsenger-scales-bot/util"
)

type Server struct {
	cfg  *config.Server
	root handler.RootHandler
}

func NewServer(cfg *config.Server) *Server {
	return &Server{
		cfg: cfg,
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
        // TODO: sentry
		// app.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
		// app.Logger.Printf("Sentry initialized")
	}
    app.Validator = util.NewDefaultValidator()    

    app.GET("/", s.root.Handle)

    addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
    app.Logger.Fatal(app.Start(addr))
}

