package db

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"gohive/internal/pb/def"
)

const (
	ITEM_ID_OFFSET = 10000
)

func GenItemID() (int64, error) {
	conn := pool.Get()
	defer conn.Close()

	id, err := redis.Int64(conn.Do("INCR", "item:count"))
	if err != nil {
		return 0, err
	}
	return ITEM_ID_OFFSET + id, nil
}

func LoadItem(id int64) (*def.Item, error) {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("GET", fmt.Sprintf("item:%d:idx", id))
	conn.Send("GET", fmt.Sprintf("item:%d:num", id))
	conn.Send("GET", fmt.Sprintf("item:%d:create_time", id))
	conn.Send("GET", fmt.Sprintf("item:%d:update_time", id))
	res, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	idx, err := redis.Int(res[0], nil)
	if err != nil {
		return nil, err
	}

	num, err := redis.Int64(res[1], nil)
	if err != nil {
		return nil, err
	}

	createTime, err := redis.Int64(res[2], nil)
	if err != nil {
		return nil, err
	}

	updateTime, err := redis.Int64(res[3], nil)
	if err != nil {
		return nil, err
	}

	item := &def.Item{
		Id:         proto.Int64(id),
		Idx:        proto.Int32(int32(idx)),
		Num:        proto.Int64(num),
		CreateTime: proto.Int64(createTime),
		UpdateTime: proto.Int64(updateTime),
	}
	return item, nil
}
