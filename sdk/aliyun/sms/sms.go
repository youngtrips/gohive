package sms

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"gohive/sdk/aliyun/config"
	"gohive/sdk/aliyun/util"
)

// {"Message":"OK","RequestId":"98E41FE5-AE8B-4FCA-9341-AD7812243B2E","BizId":"108998076955^1112054108617","Code":"OK"}

type respInfo struct {
	Message   string `json:"Message"`
	RequestId string `json:"RequestId"`
	BizId     string `json:"BizId"`
	Code      string `json:"Code"`
}

func SendCode(phoneNumbser string, code string) error {

	smsCfg := config.SDKParam.SMS

	params := make([]util.ParamInfo, 0)
	params = append(params, util.ParamInfo{"Action", "SendSms"})
	params = append(params, util.ParamInfo{"Version", "2017-05-25"})
	params = append(params, util.ParamInfo{"RegionId", "cn-hangzhou"})
	params = append(params, util.ParamInfo{"PhoneNumbers", phoneNumbser})
	params = append(params, util.ParamInfo{"SignName", smsCfg.SignName})
	params = append(params, util.ParamInfo{"TemplateParam", "{\"code\":\"" + code + "\"}"})
	params = append(params, util.ParamInfo{"TemplateCode", smsCfg.TemplateCode})
	params = append(params, util.ParamInfo{"OutId", "123"})

	url := util.BuildParam(config.SDKParam.AccessKey, config.SDKParam.AccessSecret, "GET", params)
	fullurl := smsCfg.Url + "/?" + url

	log.Info("fullurl: ", fullurl)
	resp, err := http.Get(fullurl)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	rInfo := &respInfo{}
	if err := json.Unmarshal([]byte(body), rInfo); err != nil {
		log.Error(err)
		return err
	}

	log.Info(string(body))

	if rInfo.Code != "OK" {
		log.Error(string(body))
		return errors.New(rInfo.Message)
	}
	return nil
}

/*
paras := make([]ParamInfo, 0)
	paras = append(paras, ParamInfo{"Action", "SendSms"})
	paras = append(paras, ParamInfo{"Version", "2017-05-25"})
	paras = append(paras, ParamInfo{"RegionId", "cn-hangzhou"})
	paras = append(paras, ParamInfo{"PhoneNumbers", "15928010910"})
	paras = append(paras, ParamInfo{"SignName", "电玩城101"})
	paras = append(paras, ParamInfo{"TemplateParam", "{\"code\":\"91753\"}"})
	paras = append(paras, ParamInfo{"TemplateCode", "SMS_78710086"})
	paras = append(paras, ParamInfo{"OutId", "123"})

	s := BuidParam("LTAIKabOLX7iW7bt", "35uo7vaJJBp7HmFbsmdLAhNI8LficI", "GET", paras)

	url := "http://dysmsapi.aliyuncs.com/?" + s
	t.Logf("%s", url)

	resp, err := http.Get(url)
	if err != nil {
		t.Logf("error: %s", err)
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Logf("error: %s", err)
		} else {
			t.Logf("resp: %s", string(body))
		}
	}
*/
