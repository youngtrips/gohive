package handler

import (
	"gohive/internal/pb/msg"
)

func OnGetRankReq(req *msg.GetRank_Req) *msg.GetRank_Res {
	res := &msg.GetRank_Res{}
	return res
}

func OnGetRankRes(res *msg.GetRank_Res) {
}
