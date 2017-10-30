package handler

import (
	"gohive/internal/pb/msg"
)

func OnSetEventReq(req *msg.SetEvent_Req) *msg.SetEvent_Res {
	res := &msg.SetEvent_Res{}
	return res
}

func OnSetEventRes(res *msg.SetEvent_Res) {
}
