package gate

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/youngtrips/anet"
	"gohive/internal/config"
	"gohive/internal/misc"
	"gohive/internal/naming"
	"gohive/internal/pb/encoding"
	pb "gohive/internal/pb/service"
	"gohive/server/gate/entity"
	"gohive/server/gate/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	MAXN_PENDING_PACKETS = 65535
	MAXN_ANET_EVENTS     = 65535
)

type RpcEvent struct {
	sa  *entity.ServerAgent
	pkt *pb.Packet
}

type Server struct {
	events         chan anet.Event
	rbuf           chan RpcEvent
	sc             chan os.Signal
	maxEvents      int32
	maxPendingPkts int32
	cfg            *config.ServerInfo
	tcpSrv         *anet.Server
	rpcSrv         *grpc.Server
}

func newRPCServer(cfg *config.GRpcServerInfo, rbuf chan RpcEvent) (*grpc.Server, error) {
	var rpc_srv *grpc.Server
	if cfg.SslEnable {
		certificate, err := tls.LoadX509KeyPair(
			cfg.SslKey.CertFile,
			cfg.SslKey.KeyFile,
		)
		certPool := x509.NewCertPool()
		bs, err := ioutil.ReadFile(cfg.SslKey.CaFile)
		if err != nil {
			return nil, err
		}
		if !certPool.AppendCertsFromPEM(bs) {
			log.Fatal("failed to append client certs")
		}
		tlsConfig := &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		}
		opt := grpc.Creds(credentials.NewTLS(tlsConfig))
		rpc_srv = grpc.NewServer(opt)
	} else {
		rpc_srv = grpc.NewServer()
	}
	return rpc_srv, nil
}

func NewServer(sc chan os.Signal, cfg *config.ServerInfo) (*Server, error) {
	maxPendingPkts := MAXN_PENDING_PACKETS
	maxEvents := cfg.ANet.MaxEvents
	if maxEvents <= 0 {
		maxEvents = MAXN_ANET_EVENTS
	}

	events := make(chan anet.Event, maxEvents)
	rbuf := make(chan RpcEvent, maxPendingPkts)
	rpcSrv, err := newRPCServer(&cfg.GRpc, rbuf)
	if err != nil {
		return nil, err
	}
	s := &Server{
		events:         events,
		rbuf:           rbuf,
		cfg:            cfg,
		sc:             sc,
		maxEvents:      int32(maxEvents),
		maxPendingPkts: int32(maxPendingPkts),
		tcpSrv:         anet.NewServer("tcp", cfg.ANet.Addr, &encoding.Protocol{}, events),
		rpcSrv:         rpcSrv,
	}
	return s, nil
}

func (s *Server) Start() error {
	// start tcp
	if err := s.tcpSrv.ListenAndServe(); err != nil {
		return err
	}

	// start rpc
	ln, err := net.Listen("tcp", s.cfg.GRpc.Addr)
	if err != nil {
		return err
	}
	pb.RegisterGameServiceServer(s.rpcSrv, s)
	go func() {
		if err := s.rpcSrv.Serve(ln); err != nil {
			log.Fatal(err)
		}
	}()

	if err := s.registry(); err != nil {
		log.Fatal(err)
	}

	log.Info("gateserver start ok...")
	loop := true
	for loop {
		select {
		case <-s.sc:
			loop = false
			s.tcpSrv.Close()
			s.rpcSrv.GracefulStop()
			log.Info("exit...")
			break
		case ev, ok := <-s.rbuf: // intra
			if ok {
				s.onPacket(ev.sa, ev.pkt)
			}
			break
		case ev, ok := <-s.events: // extra
			if ok {
				log.Info(ev)
				s.onEvent(ev)
			}
			break
		}
	}
	return nil
}

func (s *Server) Stop() {
	close(s.events)
	close(s.rbuf)
}

func (s *Server) onPacket(sa *entity.ServerAgent, pkt *pb.Packet) {
	handler.IntraDispatch(sa, pkt)
}

func (s *Server) onEvent(ev anet.Event) {
	switch ev.Type {
	case anet.EVENT_ACCEPT:
		ua := entity.NewUserAgent(ev.Session, s.maxEvents, s)
		if ua != nil {
			go ua.Start()
		} else {
			ev.Session.Close()
		}
		break
	}
}

func (s *Server) OnMessage(ua *entity.UserAgent, api string, payload interface{}) {
	log.Info("onMessage: ", api)
	handler.Dispatch(ua, api, payload)
}

func (s *Server) recv(stream pb.GameService_StreamServer, sess_die chan struct{}) chan *pb.Packet {
	ch := make(chan *pb.Packet, 1)
	go func() {
		defer func() {
			close(ch)
		}()
		for {
			in, err := stream.Recv()
			if err == io.EOF { // client closed
				return
			}

			if err != nil {
				log.Error(err)
				return
			}
			select {
			case ch <- in:
			case <-sess_die:
			}
			log.Info("stream recv loop...")
		}
	}()
	return ch
}

func (self *Server) Stream(stream pb.GameService_StreamServer) error {

	log.Info("open tunnel...")
	md, ok := metadata.FromIncomingContext(stream.Context())
	// read metadata from context
	if !ok {
		log.Error("cannot read metadata from context")
		return errors.New("cannot read metadata from context")
	}
	// read key
	if len(md["server"]) == 0 {
		log.Error("cannot read key:server_id from metadata")
		return errors.New("cannot read key:server id from metadata")
	}

	// parse server_id
	serverId, err := strconv.Atoi(md["server"][0])
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("server: ", serverId)

	sess_die := make(chan struct{})

	defer func() {
		close(sess_die)
	}()

	packets_i := self.recv(stream, sess_die)
	packets_o := make(chan *pb.Packet, MAXN_PENDING_PACKETS)
	ent := entity.NewServerAgent(int32(serverId), packets_o)
	defer func() {
		close(packets_o)
	}()

	for {
		select {
		case p, ok := <-packets_i:
			if ok {
				self.rbuf <- RpcEvent{ent, p}
			} else {
				log.Error("stream closed...")
				return nil
			}
			break
		case p, ok := <-packets_o:
			if ok {
				if err := stream.Send(p); err != nil {
					log.Error("send: ", err)
					return err
				}
			}
		}
	}
	return nil
}

func (s *Server) registry() error {
	log.Info("privateKey: ", entity.PRIVATE_KEY)
	log.Info("pubKey: ", entity.PUBLIC_KEY)

	p, err := naming.New(s.cfg.Etcd.Endpoints, nil, nil, false)
	if err != nil {
		return err
	}
	log.Info(p, err)
	go p.Run()

	wanIp, err := misc.GetWanIP()
	if err != nil {
		return err
	}

	port := strings.Split(s.cfg.GRpc.Addr, ":")[1]
	anetPort := strings.Split(s.cfg.ANet.Addr, ":")[1]
	advertiseAddr := s.cfg.ANet.AdvertiseAddr
	p.Put("gate", s.cfg.Id, wanIp+":"+port, 5000,
		naming.Param{"key", entity.PRIVATE_KEY},
		naming.Param{"ip", advertiseAddr},
		naming.Param{"anetPort", anetPort})
	return nil
}
