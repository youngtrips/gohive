package handler

import (
	"gohive/internal/pb/msg"
)

func OnGetMailsReq(req *msg.GetMails_Req) *msg.GetMails_Res {
	res := &msg.GetMails_Res{}
	return res
}

func OnGetMailsRes(res *msg.GetMails_Res) {
}
