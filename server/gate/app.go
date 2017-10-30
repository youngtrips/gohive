package gate

import (
	"os"

	"gohive/internal/config"
)

func Run(sc chan os.Signal, cfg *config.ServerInfo) {
	srv, _ := NewServer(sc, cfg)
	srv.Start()
	srv.Stop()
}
