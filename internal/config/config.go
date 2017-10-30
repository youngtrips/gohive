package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type TlsInfo struct {
	CaFile   string `json:"ca_file"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type EtcdInfo struct {
	Endpoints             string  `json:"endpoints"`
	DialTimeout           int32   `json:"dial_timeout"`
	InsecureTransport     bool    `json:"insecure_transport"`
	InsecureSkipTlsVerify bool    `json:"insecure_skip_tls_verify"`
	Security              TlsInfo `json:"security"`
}

type RedisInfo struct {
	Host string `json:"host"`
	Port int32  `json:"port"`
}

type SslKeyInfo struct {
	CaFile   string `json:"ca_file"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type GRpcServerInfo struct {
	Addr      string     `json:"addr"`
	SslEnable bool       `json:"ssl_enable"`
	SslKey    SslKeyInfo `json:"ssl_key"`
}

type BSScopeInfo struct {
	Scope    int32  `json:"scope"`
	Priority int32  `json:"priority"`
	BS       string `json:"bs"`
}

type BSInfo struct {
	Default  string              `json:"default`
	Services map[string][]string `json:"services"`
	Scopes   []*BSScopeInfo      `json:"scopes"`
}

type ANetInfo struct {
	Addr          string `json:"addr"`
	AdvertiseAddr string `json:"advertise_addr"`
	MaxEvents     int32  `json:"max_events"`
}

type HttpInfo struct {
	Mode string `json:"mode"`
	Host string `json:"host"`
	Port int32  `json:"port"`
}

type ServerInfo struct {
	Id    int32          `json:"id"`
	ANet  ANetInfo       `json:"anet"`
	GRpc  GRpcServerInfo `json:"grpc"`
	Http  HttpInfo       `json:"http"`
	Redis RedisInfo      `json:"redis"`
	Etcd  EtcdInfo       `json:"etcd"`
}

var (
	PWD string
)

func init() {
	PWD, _ = os.Getwd()
}

func Load(cfgfile string) (*ServerInfo, error) {
	data, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return nil, err
	}

	c := &ServerInfo{}
	if err := json.Unmarshal([]byte(data), c); err != nil {
		return nil, err
	}
	return c, nil
}
