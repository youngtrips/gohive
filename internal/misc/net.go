package misc

import (
	"net"
	"strings"
)

func GetWanIP() (string, error) {
	conn, err := net.Dial("udp", "www.baidu.com:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}
