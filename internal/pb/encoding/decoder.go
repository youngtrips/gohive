package encoding

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
)

func Decode(typeName string, payload []byte) (proto.Message, error) {
	//fmt.Printf("typeName: [%s]\n", typeName)
	typeValue := proto.MessageType(typeName)
	if typeValue == nil {
		return nil, errors.New(fmt.Sprintf("no such protocal type: %s", typeName))
	}

	typeValue = typeValue.Elem()

	value := reflect.New(typeValue)
	msg := value.Interface().(proto.Message)
	if err := proto.Unmarshal(payload, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
