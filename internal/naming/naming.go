package naming

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	etcdclient "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc"
)

const (
	MAXN_TIMEOUT        = 100 * time.Millisecond
	VERSION             = "v1"
	SERVICE_PREFIX      = "/backends"
	DEFAULT_DIALTIMEOUT = 2 * time.Second
	FPS                 = 10
)

type Param struct {
	Key string
	Val string
}

type registryInfo struct {
	Addr   string            `json:"addr"`
	Params map[string]string `json:"params"`
}

type ConnInfo struct {
	Id     int32
	Addr   string
	Params map[string]string
	Conn   *grpc.ClientConn
}

type serviceInfo struct {
	name        string
	next        int
	conns       []*ConnInfo
	connMapping map[int32]*ConnInfo // id ==> info
	sync.RWMutex
}

type Pool struct {
	sync.Mutex
	clt      *etcdclient.Client
	leases   map[string]etcdclient.LeaseID
	ctx      context.Context
	cancel   context.CancelFunc
	services map[string]*serviceInfo
	names    map[string]bool
	notify   chan *ConnInfo
	do_conn  bool
}

func New(endpoints string, tls *tls.Config, notify chan *ConnInfo, do_conn bool, names ...string) (*Pool, error) {
	addrs := strings.Split(endpoints, ",")
	clt, err := etcdclient.New(etcdclient.Config{
		Endpoints:   addrs,
		DialTimeout: DEFAULT_DIALTIMEOUT,
		TLS:         tls,
	})
	if err != nil {
		return nil, err
	}

	baseCtx := context.TODO()
	ctx, cancel := context.WithCancel(baseCtx)

	p := &Pool{
		clt:      clt,
		leases:   make(map[string]etcdclient.LeaseID),
		ctx:      ctx,
		cancel:   cancel,
		services: make(map[string]*serviceInfo),
		names:    make(map[string]bool),
		notify:   notify,
		do_conn:  do_conn,
	}
	for _, n := range names {
		p.names[n] = true
	}
	p.init()
	go p.watch()
	return p, nil
}

func (p *Pool) init() {
	key := path.Join(SERVICE_PREFIX, VERSION)
	if resp, err := p.clt.Get(context.Background(), key, etcdclient.WithPrefix()); err != nil {
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
			if !p.names[fields[3]] {
				continue
			}
			id, err := strconv.Atoi(fields[4])
			if err != nil {
				log.Warn("invalid service id: ", fields[4], " ", err)
				continue
			}

			key := path.Join("/", fields[1], fields[2], fields[3])
			p.add(key, int32(id), string(ev.Value))
		}
	}
}

func (p *Pool) put(name string, id int32, val string) error {
	key := path.Join(SERVICE_PREFIX, VERSION, name, fmt.Sprintf("%d", id))
	_, err := p.clt.Put(context.Background(), key, val)
	return err
}

func (p *Pool) putWithTTL(name string, id int32, val string, ttl int64) error {
	resp, err := p.clt.Grant(context.TODO(), ttl)
	if err != nil {
		return err
	}

	key := path.Join(SERVICE_PREFIX, VERSION, name, fmt.Sprintf("%d", id))
	_, err = p.clt.Put(context.TODO(), key, val, etcdclient.WithLease(resp.ID))
	if err != nil {
		return err
	}

	p.leases[name] = resp.ID
	return nil
}

func (p *Pool) Put(name string, id int32, addr string, ttl int64, params ...Param) error {
	///
	info := &registryInfo{
		Addr:   addr,
		Params: make(map[string]string),
	}
	for _, p := range params {
		info.Params[p.Key] = p.Val
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	///
	log.Info(name, "/", id, addr, ", ", ttl)
	if ttl > 0 {
		return p.putWithTTL(name, id, string(data), ttl)
	}
	return p.put(name, id, string(data))
}

func (p *Pool) Get(name string) (string, error) {
	return "", nil
}

func (p *Pool) Run() {
	interval := time.Duration(1000 / FPS)
	tick := time.NewTimer(0)
	for {
		select {
		case _ = <-tick.C:
			p.onTick()
			tick.Reset(interval * time.Millisecond)
			break
		}
	}
}

func (p *Pool) Close() {
	p.cancel()
	p.clt.Close()
}

func (p *Pool) onTick() {
	pending := make([]string, 0)
	for key, id := range p.leases {
		// to renew the lease only once
		_, kaerr := p.clt.KeepAliveOnce(context.TODO(), id)
		if kaerr != nil {
			log.Error(kaerr)
			pending = append(pending, key)
		}
	}
	for _, key := range pending {
		delete(p.leases, key)
	}
}

func (p *Pool) watch() {
	key := path.Join(SERVICE_PREFIX, VERSION)
	rch := p.clt.Watch(context.Background(), key, etcdclient.WithPrefix())
	running := true
	for running {
		select {
		case wresp := <-rch:
			p.onUpdate(wresp)
			break
		case <-p.ctx.Done():
			running = false
			break
		}
	}

}
func (p *Pool) onUpdate(wresp etcdclient.WatchResponse) {
	for _, ev := range wresp.Events {
		switch ev.Type {
		case mvccpb.PUT:
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
			p.add(key, int32(id), string(ev.Kv.Value))

			break
		case mvccpb.DELETE:
			log.Info("DEL: ", string(ev.Kv.Key), ", ", string(ev.Kv.Value))
			break
		}
	}
}

func (p *Pool) add(name string, id int32, val string) {
	go func() {
		info := &registryInfo{}
		if err := json.Unmarshal([]byte(val), info); err != nil {
			log.Error("invald registryInfo: ", err)
			return
		}

		var conn *grpc.ClientConn = nil
		var err error
		if p.do_conn {
			if conn, err = grpc.Dial(info.Addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(MAXN_TIMEOUT)); err != nil {
				log.Errorf("open connect to service falied: %s/%d/%s %s", name, id, info.Addr, err)
				return
			}
		}

		p.Lock()
		defer p.Unlock()
		sInfo, present := p.services[name]
		if !present {
			sInfo = &serviceInfo{
				name:        name,
				next:        0,
				conns:       make([]*ConnInfo, 0),
				connMapping: make(map[int32]*ConnInfo),
			}
			p.services[name] = sInfo
		}

		cInfo, present := sInfo.connMapping[id]
		if !present {
			cInfo = &ConnInfo{
				Id:     id,
				Addr:   info.Addr,
				Params: info.Params,
				Conn:   conn,
			}
			sInfo.conns = append(sInfo.conns, cInfo)
			sInfo.connMapping[id] = cInfo
		} else {
			if cInfo.Conn != nil {
				cInfo.Conn.Close()
			}
			cInfo.Addr = info.Addr
			cInfo.Params = info.Params
			cInfo.Conn = conn
		}
		log.Infof("add service: %s %d %s", name, id, info.Addr)
		if p.notify != nil {
			p.notify <- cInfo
		}
	}()
}
