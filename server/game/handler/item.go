package handler

import (
	"gohive/internal/pb/msg"
)

func OnGiftItemsReq(req *msg.GiftItems_Req) *msg.GiftItems_Res {
	res := &msg.GiftItems_Res{}
	return res
}

func OnGiftItemsRes(res *msg.GiftItems_Res) {
}

func OnGetItemsReq(req *msg.GetItems_Req) *msg.GetItems_Res {
	res := &msg.GetItems_Res{}
	return res
}

func OnGetItemsRes(res *msg.GetItems_Res) {
}
