package encoding

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/golang/protobuf/proto"
)

type Protocol struct {
}

func (self *Protocol) Encode(api string, msg interface{}) ([]byte, error) {
	//log.Printf("encode: %s", api)

	// name_size + name + body
	bufLen := new(bytes.Buffer)
	apiLen := uint16(len(api))

	if err := binary.Write(bufLen, binary.LittleEndian, apiLen); err != nil {
		return nil, err
	}

	///
	body := new(bytes.Buffer)

	if err := binary.Write(body, binary.LittleEndian, []byte(api)); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err := binary.Write(body, binary.LittleEndian, data); err != nil {
		return nil, err
	}

	b := body.Bytes()

	buf := make([]byte, 0)
	buf = append(buf, bufLen.Bytes()...)
	buf = append(buf, b...)
	return buf, nil
}

func (self *Protocol) Decode(data []byte) (string, interface{}, error) {
	var size int16

	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.LittleEndian, &size); err != nil {
		log.Printf("get apiname size failed: %s", err)
		return "", nil, err
	}

	log.Printf("api size: %d", size)

	api := make([]byte, size)
	if err := binary.Read(buf, binary.LittleEndian, api); err != nil {
		log.Printf("get api failed: %s", err)
		return "", nil, err
	}

	msg, err := Decode(string(api), data[size+2:])
	return string(api), msg, err
}
