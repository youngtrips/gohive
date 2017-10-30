package emulator

import (
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/youngtrips/anet"
	"gohive/internal/misc"
	"gohive/internal/pb/encoding"
	"gohive/internal/pb/msg"
)

const (
	FPS = 100
)

type Agent struct {
	session      *anet.Session
	events       chan anet.Event
	tick         int64
	serverId     int32
	running      bool
	connected    bool
	loginTime    int64
	lastMoveTime int64
	addr         string
	lastSync     int64
	token        string
}

func NewAgent(token string) *Agent {
	self := &Agent{
		token:        token,
		session:      nil,
		events:       make(chan anet.Event, 1024),
		tick:         0,
		serverId:     0,
		running:      true,
		connected:    false,
		loginTime:    0,
		lastMoveTime: 0,
		lastSync:     0,
	}
	return self
}

func (self *Agent) Start(addr string) {
	self.addr = addr
	self.session = anet.ConnectTo("tcp4", addr, &encoding.Protocol{}, self.events, false)
	log.Printf("start connect to %s...", addr)
	go self.loop()
}

func (self *Agent) Send(msg proto.Message) {
	if self.session != nil {
		self.session.Send(proto.MessageName(msg), msg)
	}
}

func (self *Agent) onEvent(ev anet.Event) {
	switch ev.Type {
	case anet.EVENT_CONNECT_SUCCESS:
		self.onConnect(ev)
		break
	case anet.EVENT_MESSAGE:
		msg := ev.Data.(*anet.Message)
		//log.Info("onMesage...", msg)
		self.onMessage(msg)
		break
	case anet.EVENT_DISCONNECT:
		self.onDisconnect()
		break
	}
}

func (self *Agent) onConnect(ev anet.Event) {
	log.Print("connect success...")
	self.session = ev.Session
	self.connected = true
	self.session.Start(nil)
	self.Login()
}

func (self *Agent) onDisconnect() {
	self.connected = false
	self.running = false
}

func (self *Agent) loop() {
	last := misc.NowMS()
	interval := int64(1000 / FPS)
	for self.running {
		self.tick++
		curr := misc.NowMS()
		dt := curr - last
		if dt >= interval {
			self.onTick(dt)
		}
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

func (self *Agent) onTick(dt int64) {
	self.proc()
}

func (self *Agent) proc() {
	select {
	case ev, ok := <-self.events:
		if ok {
			log.Print("event: ", ev)
			self.onEvent(ev)
		}
	default:
		break
	}
}

func (self *Agent) Login() {
	req := &msg.Login_Req{
		Token: proto.String(self.token),
	}
	self.Send(req)
}
