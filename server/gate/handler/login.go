package handler

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gohive/internal/pb/def"
	"gohive/internal/pb/msg"
	"gohive/server/db"
	"gohive/server/gate/entity"
)

//anonymous
func onLoginReq(ua *entity.UserAgent, req *msg.Login_Req) {
	log.Info("login token: ", req.GetToken())
	// check token
	accId, code := db.CheckToken(req.GetToken(), entity.PUBLIC_KEY)
	if code != int32(def.RC_OK) {
		log.Info("login code: ", code)
		res := &msg.Login_Res{
			Code: proto.Int32(code),
		}
		ua.Send(res)
		return
	}
	ua.BindAccount(accId)
	log.Info("AccountID: ", accId)
	sa := entity.GetServerAgent(ua.Server)
	if sa == nil {
		sa = entity.RandServerAgent()
	}
	if sa == nil {
		res := &msg.Login_Res{
			Code: proto.Int32(int32(def.RC_GS_NOT_CONNECTED)),
		}
		ua.Send(res)
	} else {
		req.Token = proto.String(fmt.Sprintf("%d", accId))
		sa.Forward(ua.Id, req)
	}
}

func onLoginRes(ua *entity.UserAgent, res *msg.Login_Res) {
}
