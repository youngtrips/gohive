package game

import (
	"os"

	"gohive/internal/config"
	"gohive/server/db"
)

func Run(sc chan os.Signal, cfg *config.ServerInfo) {
	db.Open(cfg.Redis.Host, cfg.Redis.Port)
	defer db.Close()

	srv := newServer(cfg)
	srv.Run(sc)
}
