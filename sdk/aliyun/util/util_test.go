package util

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestBuildParam(t *testing.T) {
	paras := make([]ParamInfo, 0)
	paras = append(paras, ParamInfo{"Action", "SendSms"})
	paras = append(paras, ParamInfo{"Version", "2017-05-25"})
	paras = append(paras, ParamInfo{"RegionId", "cn-hangzhou"})
	paras = append(paras, ParamInfo{"PhoneNumbers", "15928010910"})
	paras = append(paras, ParamInfo{"SignName", "电玩城101"})
	paras = append(paras, ParamInfo{"TemplateParam", "{\"code\":\"91753\"}"})
	paras = append(paras, ParamInfo{"TemplateCode", "SMS_78710086"})
	paras = append(paras, ParamInfo{"OutId", "123"})

	s := BuildParam("LTAIKabOLX7iW7bt", "35uo7vaJJBp7HmFbsmdLAhNI8LficI", "GET", paras)

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
}
