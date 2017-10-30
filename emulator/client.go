package emulator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/youngtrips/anet"
	"gohive/internal/pb/def"
	"gohive/internal/pb/msg"
)

const (
	BASE_URL = "http://qeetap.com:8082"
	//BASE_URL = "http://127.0.0.1:8082"
)

type gateconf struct {
	id   int
	name string
	ip   string
	port string
}

var _ = proto.Int32

func pp(format string, v ...interface{}) {
	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Printf(format, v...)
}

func init() {
}

type Client struct {
	agent  *Agent
	gate   gateconf
	events chan anet.Event
}

func NewClient() *Client {
	return &Client{
		events: make(chan anet.Event, 65535),
	}
}

func (self *Client) Do_registry(username string, password string) {
	const REG_URL = BASE_URL + "/v1/account/registry"

	v := url.Values{}
	v.Set("username", username)
	v.Set("password", password)
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))

	client := &http.Client{}
	req, _ := http.NewRequest("POST", REG_URL, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		pp("auth failed: %s", err)
		return
	}
	res := &msg.Auth_Res{}
	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		pp("auth failed: %s", err)
		return
	} else {
		pp("body: %s", string(body))
		if err := json.Unmarshal([]byte(body), res); err != nil {
			pp("auth failed: %s", err)
			return
		}
		if res.GetCode() != int32(def.RC_OK) {
			pp("auth failed: code=%d", res.Code)
			return
		}
	}
	pp(res.GetToken())
	pp(res.GetAddr())
	if res.GetToken() != "" && res.GetAddr() != "" {
		self.agent = NewAgent(res.GetToken())
		self.agent.Start(res.GetAddr())
	}
}

func (self *Client) Do_login(username string, password string) {
	const AUTH_URL = BASE_URL + "/v1/account/auth"

	url := fmt.Sprintf("%s?username=%s&password=%s", AUTH_URL, username, password)
	resp, err := http.Get(url)
	if err != nil {
		pp("auth failed: %s", err)
		return
	}
	res := &msg.Auth_Res{}
	defer resp.Body.Close()
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		pp("auth failed: %s", err)
		return
	} else {
		pp("body: %s", string(body))
		if err := json.Unmarshal([]byte(body), res); err != nil {
			pp("auth failed: %s", err)
			return
		}
		if res.GetCode() != int32(def.RC_OK) {
			pp("auth failed: code=%d", res.Code)
			return
		}
	}
	pp(res.GetToken())
	pp(res.GetAddr())
	if res.GetToken() != "" && res.GetAddr() != "" {
		self.agent = NewAgent(res.GetToken())
		self.agent.Start(res.GetAddr())
	}
}
