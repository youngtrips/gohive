package handler

import (
	"gohive/internal/pb/msg"
)

func OnGetUserReq(req *msg.GetUser_Req) *msg.GetUser_Res {
	res := &msg.GetUser_Res{}
	return res
}

func OnGetUserRes(res *msg.GetUser_Res) {
}
