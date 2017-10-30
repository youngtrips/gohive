package emulator

import (
	"github.com/youngtrips/anet"
	"gohive/internal/pb/def"
	"gohive/internal/pb/msg"
)

func (self *Agent) onMessage(m *anet.Message) {
	switch m.Api {
	case "msg.Login.Res":
		onLoginRes(self, m.Payload.(*msg.Login_Res))
		break
	}
}

func onLoginRes(agent *Agent, res *msg.Login_Res) {
	pp("loginRes code: %d", res.GetCode())
	if res.GetCode() == int32(def.RC_OK) {
		pp("user: %+v", res.User)
	}
}
