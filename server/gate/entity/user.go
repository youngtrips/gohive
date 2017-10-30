package entity

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/youngtrips/anet"
	"gohive/internal/pb/msg"
)

type Dispatcher interface {
	OnMessage(ua *UserAgent, api string, payload interface{})
}

type UserAgent struct {
	Id         int64
	Account    int64
	UserId     int64
	Server     int32
	dispatcher Dispatcher
	session    *anet.Session
	events     chan anet.Event
	signal     chan int32
	running    bool
	sync.Mutex
}

var (
	_uas      map[int64]*UserAgent
	_uas_lock sync.Mutex
)

func init() {
	_uas = make(map[int64]*UserAgent)
}

func (ua *UserAgent) Send(msg proto.Message) {
	if ua.session != nil {
		ua.session.Send(proto.MessageName(msg), msg)
	}
}

func (ua *UserAgent) RawSend(api string, data []byte) {
	if ua.session != nil {
		ua.session.RawSend(api, data)
	}
}

func GetUserAgent(id int64) *UserAgent {
	_uas_lock.Lock()
	defer _uas_lock.Unlock()
	return _uas[id]
}

func NewUserAgent(sess *anet.Session, maxEvents int32, dispatcher Dispatcher) *UserAgent {
	if sess == nil {
		return nil
	}

	ua := &UserAgent{
		Id:         0,
		UserId:     0,
		Server:     0,
		session:    sess,
		dispatcher: dispatcher,
		events:     make(chan anet.Event, maxEvents),
		signal:     make(chan int32),
		running:    true,
	}

	if sess != nil {
		ua.Id = sess.ID()
	}

	_uas_lock.Lock()
	defer _uas_lock.Unlock()
	_uas[ua.Id] = ua
	return ua
}

func (ua *UserAgent) Start() {
	defer func() {
		if ua.session != nil {
			ua.session.Close()
		}
		close(ua.events)
		ua.session = nil
	}()
	ua.session.Start(ua.events)
	for ua.running {
		select {
		case <-ua.signal:
			ua.running = false
			break
		case ev, ok := <-ua.events:
			if ok {
				ua.onEvent(ev)
			}
			break
		}
	}
}

func (ua *UserAgent) onEvent(ev anet.Event) {
	switch ev.Type {
	case anet.EVENT_DISCONNECT:
		log.Info("connect is closed by remote peer: ", ua.Id)
		ua.onLogout()
		ua.signal <- 1
		break
	case anet.EVENT_MESSAGE:
		msg := ev.Data.(*anet.Message)
		if ua.dispatcher != nil && msg != nil {
			ua.dispatcher.OnMessage(ua, msg.Api, msg.Payload)
		}
		break
	}
}

func (ua *UserAgent) onLogout() {
	log.Info("logout, server: ", ua.Server, " ", ua.UserId)
	if ua.Server > 0 && ua.Account > 0 {
		if sa := GetServerAgent(ua.Server); sa != nil {
			req := &msg.Logout_Req{
				Account: proto.Int64(ua.Account),
			}
			sa.Forward(ua.Id, req)
		}
	}
}

func (ua *UserAgent) BindAccount(accId int64) {
	ua.Lock()
	defer ua.Unlock()
	ua.Account = accId
}
