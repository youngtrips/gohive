package handler

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gohive/internal/pb/def"
	"gohive/internal/pb/msg"
	"gohive/server/game/entity"
)

func OnLoginReq(req *msg.Login_Req) *msg.Login_Res {
	log.Infof("onLogin: %+v", req)

	accId, _ := strconv.ParseInt(req.GetToken(), 10, 64)
	user := entity.LoadUser(accId)
	if user == nil {
		return &msg.Login_Res{
			Code: proto.Int32(int32(def.RC_LOGIN_USER_CREATE_FAILED)),
		}
	}
	res := &msg.Login_Res{
		Code: proto.Int32(int32(def.RC_OK)),
		User: user.PB,
	}
	for _, item := range user.Items {
		res.User.Items = append(res.User.Items, item)
	}
	return res
}

func OnLoginRes(res *msg.Login_Res) {
}
