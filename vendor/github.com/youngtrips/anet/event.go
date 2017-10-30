package anet

const (
	EVENT_ACCEPT          = 0x01
	EVENT_CONNECT_SUCCESS = 0x02
	EVENT_CONNECT_FAILED  = 0x04
	EVENT_DISCONNECT      = 0x08
	EVENT_MESSAGE         = 0x10
	EVENT_RECV_ERROR      = 0x20
	EVENT_SEND_ERROR      = 0x40
)

type Event struct {
	Type    int8
	Session *Session
	Data    interface{}
}

func newEvent(typ int8, session *Session, data interface{}) Event {
	return Event{
		Type:    typ,
		Session: session,
		Data:    data,
	}
}
