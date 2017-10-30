package handler

import (
	"gohive/internal/pb/msg"
)

func OnEchoReq(req *msg.Echo_Req) *msg.Echo_Res {
	res := &msg.Echo_Res{}
	return res
}

func OnEchoRes(res *msg.Echo_Res) {
}
