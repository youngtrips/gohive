package handler

import (
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"gohive/internal/pb/def"
	"gohive/internal/pb/msg"
	"gohive/internal/pb/service"
	"gohive/server/gate/entity"
)

func Dispatch(ua *entity.UserAgent, api string, m interface{}) {
	log.Info("dispatch: ", api, m)
	if api == "msg.Login.Req" {
		onLoginReq(ua, m.(*msg.Login_Req))
	} else {
		if ua.Server > 0 {
			forward(ua, api, m)
		} else {
			res := &msg.Login_Res{
				Code: proto.Int32(int32(def.RC_LOGIN_TOKEN_INVALID)),
			}
			ua.Send(res)
		}
	}
}

func forward(ua *entity.UserAgent, api string, m interface{}) {
	sa := entity.GetServerAgent(ua.Server)
	if sa == nil {
		log.Warn("invalid message: %s", api)
		return
	}

	data, err := proto.Marshal(m.(proto.Message))
	if err != nil {
		log.Error("invalid message: %s %s", api, err)
		return
	}

	pkt := &service.Packet{
		Peer:    ua.Id,
		Api:     api,
		Payload: data,
	}
	sa.Send(pkt)
}
