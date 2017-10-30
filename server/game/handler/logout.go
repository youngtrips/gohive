package handler

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gohive/internal/pb/def"
	"gohive/internal/pb/msg"
	"gohive/server/game/entity"
)

func OnLogoutReq(req *msg.Logout_Req) *msg.Logout_Res {
	res := &msg.Logout_Res{
		Code: proto.Int32(int32(def.RC_OK)),
	}
	entity.DelUser(req.GetAccount())

	log.Info("logout: account:", req.GetAccount())
	return res
}

func OnLogoutRes(res *msg.Logout_Res) {
}
