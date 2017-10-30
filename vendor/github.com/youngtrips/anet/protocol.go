package anet

type Protocol interface {
	Encode(api string, data interface{}) ([]byte, error)
	Decode(data []byte) (string, interface{}, error)
}
