package entity

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"gohive/internal/pb/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Event struct {
	Id  int32
	Pkt *service.Packet
}

type Session struct {
	id   int32
	rbuf chan Event
	wbuf chan *service.Packet
	md   map[string]string
	conn *grpc.ClientConn
}

func NewSession(id int32, conn *grpc.ClientConn, rbuf chan Event, md map[string]string) *Session {
	return &Session{
		id:   id,
		rbuf: rbuf,
		wbuf: make(chan *service.Packet, 65535),
		md:   md,
		conn: conn,
	}
}

func (s *Session) Start() error {
	go func() {
		//defer conn.Close()
		c := service.NewGameServiceClient(s.conn)
		for retry := 0; ; retry++ {
			ctx := metadata.NewContext(context.Background(), metadata.New(s.md))
			ctx, cancel := context.WithCancel(ctx)
			if stream, err := c.Stream(ctx); err != nil {
				time.Sleep(time.Duration(retry+1) * time.Second)
			} else {
				go func() {
					for {
						if p, err := stream.Recv(); err != nil {
							log.Error("recv: ", err)
							cancel()
							break
						} else {
							s.rbuf <- Event{s.id, p}
						}
					}
				}()
				running := true
				for running {
					select {
					case p, ok := <-s.wbuf:
						if ok {
							if err := stream.Send(p); err != nil {
								log.Error("Send: ", err)
								running = false
							}
						}
						break
					}
				}
				stream.CloseSend()
			}
		}
	}()
	return nil
}

func (s *Session) Send(pkt *service.Packet) {
	s.wbuf <- pkt
}
