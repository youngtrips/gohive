package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

var (
	pool *redis.Pool
)

func init() {
	pool = nil
}

func Open(host string, port int32) {
	addr := fmt.Sprintf("%s:%d", host, port)
	pool = &redis.Pool{
		MaxIdle:     8,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}

func get(key string, val interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	buf, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(buf), val); err != nil {
		return err
	}
	return nil
}

func set(key string, val interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	buf, err := json.Marshal(val)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, string(buf))
	return err
}

func del(key string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

//

func strKey(key interface{}) (string, error) {
	var s string
	switch key.(type) {
	case int:
		s = fmt.Sprintf("%d", key.(int))
		break
	case int8:
		s = fmt.Sprintf("%d", key.(int8))
		break
	case int16:
		s = fmt.Sprintf("%d", key.(int16))
		break
	case int32:
		s = fmt.Sprintf("%d", key.(int32))
		break
	case int64:
		s = fmt.Sprintf("%d", key.(int64))
		break
	case string:
		s = fmt.Sprintf("%s", key.(string))
		break
	default:
		return "", errors.New("invalid key type")
	}
	return s, nil
}

func genPBKey(key interface{}, dst proto.Message) (string, error) {
	s := proto.MessageName(dst)
	skey, err := strKey(key)
	if err != nil {
		return "", nil
	}
	return s + "." + skey, nil
}

func GetPB(key interface{}, dst proto.Message) error {
	fullKey, err := genPBKey(key, dst)
	if err != nil {
		return nil
	}

	conn := pool.Get()
	defer conn.Close()

	buf, err := redis.String(conn.Do("GET", fullKey))
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(buf), dst); err != nil {
		return err
	}
	return nil
}

func SetPB(key interface{}, dst proto.Message) error {
	skey, err := strKey(key)
	if err != nil {
		return err
	}
	pbName := proto.MessageName(dst)
	fullKey := pbName + "." + skey

	conn := pool.Get()
	defer conn.Close()

	buf, err := json.Marshal(dst)
	if err != nil {
		return err
	}

	if _, err := conn.Do("SET", fullKey, string(buf)); err != nil {
		return err
	}

	if _, err := conn.Do("SADD", pbName, skey); err != nil {
		return err
	}
	return nil
}

func DelPB(key interface{}, dst proto.Message) error {
	skey, err := strKey(key)
	if err != nil {
		return err
	}
	pbName := proto.MessageName(dst)
	fullKey := pbName + "." + skey

	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SREM", pbName, skey); err != nil {
		return err
	}

	if _, err := conn.Do("DEL", fullKey); err != nil {
		return err
	}
	return nil
}

func GetPBList(pbName string) ([]proto.Message, error) {

	conn := pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("SMEMBERS", pbName))
	if err != nil {
		return nil, err
	}
	pbs := make([]proto.Message, 0)
	for _, key := range keys {
		fullKey := pbName + "." + key
		buf, err := redis.String(conn.Do("GET", fullKey))
		if err != nil {
			continue
		}
		if pb, err := decode(pbName, []byte(buf)); err != nil {
			continue
		} else {
			pbs = append(pbs, pb)
		}
	}
	return pbs, nil
}

func decode(typeName string, payload []byte) (proto.Message, error) {
	typeValue := proto.MessageType(typeName)
	if typeValue == nil {
		return nil, errors.New(fmt.Sprintf("no such protocal type: %s", typeName))
	}

	typeValue = typeValue.Elem()

	value := reflect.New(typeValue)
	msg := value.Interface().(proto.Message)
	if err := json.Unmarshal(payload, msg); err != nil {
		return nil, err
	}
	return msg, nil
}
