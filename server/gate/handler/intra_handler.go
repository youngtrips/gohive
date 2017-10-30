package handler

import (
	log "github.com/Sirupsen/logrus"
	pb "gohive/internal/pb/service"
	"gohive/server/gate/entity"
)

func IntraDispatch(sa *entity.ServerAgent, pkt *pb.Packet) {
	log.Info("onPacket: ", pkt)
	ua := entity.GetUserAgent(pkt.Peer)
	if ua == nil {
		log.Warn("no found userAgent: ", pkt.Peer)
		return
	}
	ua.RawSend(pkt.Api, pkt.Payload)
	if pkt.Api == "msg.Login.Res" {
		ua.Server = sa.Id
	} else if pkt.Api == "msg.Echo.Req" {

		sa.Send(&pb.Packet{
			Api:     "msg.Echo.Res",
			Payload: pkt.Payload,
		})
	}
}
