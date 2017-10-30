package web

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gohive/internal/config"
)

func Run(cfg *config.ServerInfo) {
	e := echo.New()
	e.Debug = true
	if cfg.Http.Mode == "debug" {
		e.Debug = true
	} else if cfg.Http.Mode == "prod" {
		e.Debug = false
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	route_init(e)

	addr := fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port)
	log.Info("listen on : ", addr)

	e.Start(addr)
	log.Info("quit...")
}
