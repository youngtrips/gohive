package uc

import (
	log "github.com/Sirupsen/logrus"
	"gohive/internal/config"
	"gohive/internal/naming"
	"gohive/server/db"
	"gohive/server/uc/entity"
	"gohive/server/uc/web"
)

func Run(cfg *config.ServerInfo) {
	log.Info(*cfg)
	db.Open(cfg.Redis.Host, cfg.Redis.Port)
	defer db.Close()

	conns := make(chan *naming.ConnInfo, 65535)
	p, err := naming.New(cfg.Etcd.Endpoints, nil, conns, false, "gate")
	if err != nil {
		log.Fatal(err)
	}
	go p.Run()
	go func() {
		for {
			select {
			case info, ok := <-conns:
				if ok {
					addr := info.Params["ip"] + ":" + info.Params["anetPort"]
					log.Infof("add service: %d %s %s", info.Id, addr, info.Params["key"])
					entity.AddSerice(info.Id, addr, info.Params["key"])
				}
				break
			}
		}
	}()

	web.Run(cfg)
}
