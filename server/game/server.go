package game

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gohive/internal/config"
	"gohive/internal/naming"
	"gohive/internal/pb/encoding"
	"gohive/internal/pb/service"
	"gohive/server/game/entity"
)

const (
	MAXN_PENDING_PACKETS = 65535
)

type Server struct {
	id       int32
	cfg      *config.ServerInfo
	events   chan entity.Event
	sessions map[int32]*entity.Session
	pool     *naming.Pool
}

func newServer(cfg *config.ServerInfo) *Server {
	events := make(chan entity.Event, MAXN_PENDING_PACKETS)
	return &Server{
		id:       cfg.Id,
		cfg:      cfg,
		events:   events,
		sessions: make(map[int32]*entity.Session),
	}
}

func (s *Server) Run(sc chan os.Signal) {
	conns := make(chan *naming.ConnInfo, 65535)
	p, err := naming.New(s.cfg.Etcd.Endpoints, nil, conns, true, "gate")
	if err != nil {
		log.Fatal(err)
	}
	go p.Run()
	s.loop(sc, conns)
}

func (s *Server) loop(sc chan os.Signal, conns chan *naming.ConnInfo) {
	log.Info("server loop start...")
	md := map[string]string{"server": fmt.Sprint(s.id)}
	running := true
	for running {
		select {
		case <-sc:
			running = false
			break
		case c, ok := <-conns:
			if ok {
				sess := entity.NewSession(c.Id, c.Conn, s.events, md)
				sess.Start()
				s.sessions[c.Id] = sess
			}
			break
		case ev, ok := <-s.events:
			if ok {
				s.onEvent(ev)
			}
			break
		}
	}
}

func (s *Server) Send(id int32, peer int64, msg proto.Message) error {
	sess := s.sessions[id]
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error("encode message failed: ", err)
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			log.Error("send pkt: ", err)
		}
	}()
	pkt := &service.Packet{
		Peer:    peer,
		Api:     proto.MessageName(msg),
		Payload: data,
	}
	sess.Send(pkt)
	return nil
}

func (s *Server) onEvent(ev entity.Event) {
	pkt := ev.Pkt
	log.Info("recv pkt: ", *pkt)
	msg, err := encoding.Decode(pkt.Api, pkt.Payload)
	if err != nil {
		log.Error("invalid message: ", pkt.Api)
	} else {
		if res := Dispatch(pkt.Api, msg); res != nil {
			s.Send(ev.Id, pkt.Peer, res)
		}
	}
}
