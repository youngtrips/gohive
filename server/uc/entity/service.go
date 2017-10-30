package entity

import (
	"sync"

	log "github.com/Sirupsen/logrus"
)

var _ = log.Info

type Service struct {
	Id     int32
	Addr   string
	Key    string
	Load   int32
	Status int32
}

var (
	_srvMapping map[int32]*Service
	_srvList    []int32
	_next       int
	_srvLock    sync.Mutex
)

func init() {
	_srvMapping = make(map[int32]*Service)
	_next = -1
	_srvList = make([]int32, 0)
}

func AddSerice(id int32, addr string, key string) {
	_srvLock.Lock()
	defer _srvLock.Unlock()

	s, present := _srvMapping[id]
	if !present {
		_srvList = append(_srvList, id)
		s = &Service{
			Id:     id,
			Addr:   addr,
			Key:    key,
			Load:   0,
			Status: 0,
		}
	} else {
		s.Addr = addr
		s.Key = key
	}
	_srvMapping[id] = s
}

func GetService() *Service {
	_srvLock.Lock()
	defer _srvLock.Unlock()

	size := len(_srvList)
	if size <= 0 {
		return nil
	}
	_next = (_next + 1) % size
	return _srvMapping[_srvList[int32(_next)]]
}
