package registry

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	etcdclient "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	VERSION        = "v1"
	SERVICE_PREFIX = "/backends"
	MAXN_TIMEOUT   = 100 * time.Millisecond
)

type ConnInfo struct {
	Id   int32
	Addr string
	Conn *grpc.ClientConn
}

type serviceInfo struct {
	name        string
	next        int
	conns       []*ConnInfo
	connMapping map[int32]*ConnInfo // id ==> info
	sync.RWMutex
}

type pool struct {
	cli      *etcdclient.Client
	services map[string]*serviceInfo
	sync.RWMutex
}

var (
	_pool *pool
)

func init() {
	_pool = &pool{
		cli:      nil,
		services: make(map[string]*serviceInfo),
	}
	var err error
	_pool.cli, err = etcdclient.NewFromConfigFile("conf/etcd/etcd.yaml")
	if err != nil {
		log.Fatal(err)
	}
	_pool.init()
}

func (self *serviceInfo) getAll() []*ConnInfo {
	self.Lock()
	defer self.Unlock()

	return self.conns
}

func (self *serviceInfo) getByID(id int32) *ConnInfo {
	self.Lock()
	defer self.Unlock()

	c, present := self.connMapping[id]
	if !present {
		return nil
	}
	return c
}

func (self *serviceInfo) getOne() *ConnInfo {
	self.Lock()
	defer self.Unlock()

	curr := self.next
	self.next = (self.next + 1) % len(self.conns)
	return self.conns[curr]
}

///
func (self *pool) init() {
	key := path.Join(SERVICE_PREFIX, VERSION)
	log.Info("key: ", key)
	if resp, err := self.cli.Get(context.Background(), key, etcdclient.WithPrefix()); err != nil {
		log.Info("failed...")
		log.Error(err)
	} else {
		log.Info(resp.Kvs)
		for _, ev := range resp.Kvs {
			log.Info(string(ev.Key), " ", string(ev.Value))
			fields := strings.Split(string(ev.Key), "/")
			if len(fields) != 5 {
				log.Warn("invalid path: ", string(ev.Key))
				continue
			}
			id, err := strconv.Atoi(fields[4])
			if err != nil {
				log.Warn("invalid service id: ", fields[4], " ", err)
				continue
			}

			key := path.Join("/", fields[1], fields[2], fields[3])
			self.add(key, int32(id), string(ev.Value))
		}
	}

	go self.watch()
}

func (self *pool) watch() {
	key := path.Join(SERVICE_PREFIX, VERSION)
	rch := self.cli.Watch(context.Background(), key, etcdclient.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if ev.Type == mvccpb.PUT {
				fields := strings.Split(string(ev.Kv.Key), "/")
				if len(fields) != 5 {
					log.Warn("invalid path: ", string(ev.Kv.Key))
					continue
				}
				id, err := strconv.Atoi(fields[4])
				if err != nil {
					log.Warn("invalid service id: ", fields[4], " ", err)
					continue
				}

				key := path.Join("/", fields[1], fields[2], fields[3])
				self.add(key, int32(id), string(ev.Kv.Value))
			}
		}
	}
}

func (self *pool) add(name string, id int32, addr string) {
	go func() {
		conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(MAXN_TIMEOUT))
		if err != nil {
			log.Errorf("open connect to service falied: %s/%d/%s %s", name, id, addr, err)
			return
		}

		self.Lock()
		defer self.Unlock()
		sInfo, present := self.services[name]
		if !present {
			sInfo = &serviceInfo{
				name:        name,
				next:        0,
				conns:       make([]*ConnInfo, 0),
				connMapping: make(map[int32]*ConnInfo),
			}
			self.services[name] = sInfo
		}

		cInfo, present := sInfo.connMapping[id]
		if !present {
			cInfo = &ConnInfo{
				Id:   id,
				Addr: addr,
				Conn: conn,
			}
			sInfo.conns = append(sInfo.conns, cInfo)
			sInfo.connMapping[id] = cInfo
		} else {
			if cInfo.Conn != nil {
				cInfo.Conn.Close()
			}
			cInfo.Addr = addr
			cInfo.Conn = conn
		}
		log.Infof("add service: %s %d %s", name, id, addr)
	}()
}

func (self *pool) put(name string, id int32, addr string) error {
	self.Lock()
	defer self.Unlock()

	key := path.Join(SERVICE_PREFIX, VERSION, name, fmt.Sprintf("%d", id))
	_, err := self.cli.Put(context.Background(), key, addr)
	return err
}

func (self *pool) get(name string) *serviceInfo {
	self.Lock()
	defer self.Unlock()

	key := path.Join(SERVICE_PREFIX, VERSION, name)
	info, present := self.services[key]
	if !present {
		return nil
	}
	return info
}

func (self *pool) getAll(name string) []*ConnInfo {
	s := self.get(name)
	if s == nil {
		return nil
	}
	return s.getAll()
}

func (self *pool) getByID(name string, id int32) *ConnInfo {
	s := self.get(name)
	if s == nil {
		return nil
	}
	return s.getByID(id)
}

func (self *pool) getOne(name string) *ConnInfo {
	s := self.get(name)
	if s == nil {
		return nil
	}
	return s.getOne()
}

//
func Put(name string, id int32, addr string) error {
	return _pool.put(name, id, addr)
}

func GetAll(name string) []*ConnInfo {
	return _pool.getAll(name)
}

func Get(name string, id int32) *ConnInfo {
	return _pool.getByID(name, id)
}

func One(name string) *ConnInfo {
	return _pool.getOne(name)
}
