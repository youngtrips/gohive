package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type SMSParamInfo struct {
	Url          string `json:"url"`
	SignName     string `json:"sign_name"`
	TemplateCode string `json:"template_code"`
}

type SDKParamInfo struct {
	AccessKey    string       `json:"access_key"`
	AccessSecret string       `json:"access_secret"`
	SMS          SMSParamInfo `json:"SMS"`
}

var (
	SDKParam SDKParamInfo
)

func init() {
	PWD, _ := os.Getwd()
	data, err := ioutil.ReadFile(filepath.Join(PWD, "conf/sdk/aliyun.json"))
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal([]byte(data), &SDKParam); err != nil {
		log.Fatal(err)
	}
}
