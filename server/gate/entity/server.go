package entity

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gohive/internal/pb/service"
)

type ServerAgent struct {
	Id   int32
	wbuf chan *service.Packet
}

var (
	_sa_lock sync.Mutex
	_sas     map[int32]*ServerAgent
)

func init() {
	_sas = make(map[int32]*ServerAgent)
}

func (sa *ServerAgent) Send(pkt *service.Packet) {
}

func (sa *ServerAgent) Forward(peer int64, msg proto.Message) error {

	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	pkt := &service.Packet{
		Peer:    peer,
		Api:     proto.MessageName(msg),
		Payload: data,
	}

	sa.wbuf <- pkt
	log.Info("send pkt: ", pkt)
	return nil
}

func NewServerAgent(id int32, wbuf chan *service.Packet) *ServerAgent {

	sa := &ServerAgent{
		Id:   id,
		wbuf: wbuf,
	}

	_sa_lock.Lock()
	defer _sa_lock.Unlock()
	_sas[id] = sa
	return sa
}

func RandServerAgent() *ServerAgent {
	_sa_lock.Lock()
	defer _sa_lock.Unlock()
	for _, sa := range _sas {
		if sa != nil {
			return sa
		}
	}
	return nil
}

func GetServerAgent(id int32) *ServerAgent {
	_sa_lock.Lock()
	defer _sa_lock.Unlock()
	return _sas[id]
}
