package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

type ParamInfo struct {
	Key string
	Val string
}

type ParamList []ParamInfo

func (p ParamList) Len() int {
	return len(p)
}

func (p ParamList) Swap(i int, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ParamList) Less(i int, j int) bool {
	return p[i].Key < p[j].Key
}

func specialUrlEncode(ctx string) string {
	ctx = url.QueryEscape(ctx)
	ctx = strings.Replace(ctx, "+", "%20", -1)
	ctx = strings.Replace(ctx, "*", "%2A", -1)
	ctx = strings.Replace(ctx, "%7E", "~", -1)
	return ctx
}

func BuildParam(accessKey string, accessSecret string, method string, params []ParamInfo) string {
	vals := make([]ParamInfo, 0)
	for _, p := range params {
		if p.Key != "Signature" {
			vals = append(vals, p)
		}
	}

	ts := time.Now().UTC().Format("2006-01-02 15:04:05")

	vals = append(vals, ParamInfo{"SignatureMethod", "HMAC-SHA1"})
	vals = append(vals, ParamInfo{"SignatureNonce", fmt.Sprintf("%d", time.Now().UnixNano())})
	vals = append(vals, ParamInfo{"AccessKeyId", accessKey})
	vals = append(vals, ParamInfo{"SignatureVersion", "1.0"})
	vals = append(vals, ParamInfo{"Timestamp", ts + "Z"})
	vals = append(vals, ParamInfo{"Format", "JSON"})

	sort.Sort(ParamList(vals))

	ctx := ""
	for idx, p := range vals {
		if idx == 0 {
			ctx += fmt.Sprintf("%s=%s", specialUrlEncode(p.Key), specialUrlEncode(p.Val))
		} else {
			ctx += fmt.Sprintf("&%s=%s", specialUrlEncode(p.Key), specialUrlEncode(p.Val))
		}
	}

	signature := sign(accessSecret+"&", method+"&"+specialUrlEncode("/")+"&"+specialUrlEncode(ctx))
	fmt.Printf("signature: %s\n", signature)
	signature = specialUrlEncode(signature)
	return ctx + "&Signature=" + signature
}

func sign(key string, ctx string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(ctx))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

/*
public class SignDemo {
    public static void main(String[] args) throws Exception {
        String accessKeyId = "testId";
        String accessSecret = "testSecret";
        java.text.SimpleDateFormat df = new java.text.SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'");
        df.setTimeZone(new java.util.SimpleTimeZone(0, "GMT"));// 这里一定要设置GMT时区
        java.util.Map<String, String> paras = new java.util.HashMap<String, String>();
        // 1. 系统参数
        paras.put("SignatureMethod", "HMAC-SHA1");
        paras.put("SignatureNonce", java.util.UUID.randomUUID().toString());
        paras.put("AccessKeyId", accessKeyId);
        paras.put("SignatureVersion", "1.0");
        paras.put("Timestamp", df.format(new java.util.Date()));
        paras.put("Format", "XML");
        // 2. 业务API参数
        paras.put("Action", "SendSms");
        paras.put("Version", "2017-05-25");
        paras.put("RegionId", "cn-hangzhou");
        paras.put("PhoneNumbers", "15300000001");
        paras.put("SignName", "阿里云短信测试专用");
        paras.put("TemplateParam", "{\"customer\":\"test\"}");
        paras.put("TemplateCode", "SMS_71390007");
        paras.put("OutId", "123");
        // 3. 去除签名关键字Key
        if (paras.containsKey("Signature"))
            paras.remove("Signature");
        // 4. 参数KEY排序
        java.util.TreeMap<String, String> sortParas = new java.util.TreeMap<String, String>();
        sortParas.putAll(paras);
        // 5. 构造待签名的字符串
        java.util.Iterator<String> it = sortParas.keySet().iterator();
        StringBuilder sortQueryStringTmp = new StringBuilder();
        while (it.hasNext()) {
            String key = it.next();
            sortQueryStringTmp.append("&").append(specialUrlEncode(key)).append("=").append(specialUrlEncode(paras.get(key)));
        }
        String sortedQueryString = sortQueryStringTmp.substring(1);// 去除第一个多余的&符号
        StringBuilder stringToSign = new StringBuilder();
        stringToSign.append("GET").append("&");
        stringToSign.append(specialUrlEncode("/")).append("&");
        stringToSign.append(specialUrlEncode(sortedQueryString));

        String sign = sign(accessSecret + "&", stringToSign.toString());

        // 6. 签名最后也要做特殊URL编码
        String signature = specialUrlEncode(sign);
        System.out.println(paras.get("SignatureNonce"));
        System.out.println("\r\n=========\r\n");
        System.out.println(paras.get("Timestamp"));
        System.out.println("\r\n=========\r\n");
        System.out.println(sortedQueryString);
        System.out.println("\r\n=========\r\n");
        System.out.println(stringToSign.toString());
        System.out.println("\r\n=========\r\n");
        System.out.println(sign);
        System.out.println("\r\n=========\r\n");
        System.out.println(signature);
        System.out.println("\r\n=========\r\n");
        // 最终打印出合法GET请求的URL
        System.out.println("http://dysmsapi.aliyuncs.com/?Signature=" + signature + sortQueryStringTmp);
    }
    public static String specialUrlEncode(String value) throws Exception {
        return java.net.URLEncoder.encode(value, "UTF-8").replace("+", "%20").replace("*", "%2A").replace("%7E", "~");
    }
    public static String sign(String accessSecret, String stringToSign) throws Exception {
        javax.crypto.Mac mac = javax.crypto.Mac.getInstance("HmacSHA1");
        mac.init(new javax.crypto.spec.SecretKeySpec(accessSecret.getBytes("UTF-8"), "HmacSHA1"));
        byte[] signData = mac.doFinal(stringToSign.getBytes("UTF-8"));
        return new sun.misc.BASE64Encoder().encode(signData);
    }
}
`
*/
