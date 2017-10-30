package db

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"gohive/internal/misc"
	"gohive/internal/pb/def"
)

const (
	USER_ID_OFFSET = 10000
)

type User struct {
	PB    *def.User
	Items map[int64]*def.Item
	Mails map[int64]bool
}

func GetUserId(accId int64) (int64, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("GET", fmt.Sprintf("user:account:%d", accId)))
}

func GenUserID() (int64, error) {
	conn := pool.Get()
	defer conn.Close()

	id, err := redis.Int64(conn.Do("INCR", "user:count"))
	if err != nil {
		return 0, err
	}
	return ACCOUNT_ID_OFFSET + id, nil
}

func CreateUser(accId int64) (*User, error) {
	id, err := GenUserID()
	if err != nil {
		return nil, err
	}
	pb := &def.User{
		Id:         proto.Int64(id),
		Account:    proto.Int64(accId),
		Name:       proto.String(fmt.Sprintf("guest-%d", id)),
		Icon:       proto.Int32(1),
		Lvl:        proto.Int32(1),
		Vip:        proto.Int32(1),
		Status:     proto.Int32(1),
		CreateTime: proto.Int64(misc.NowMS()),
		UpdateTime: proto.Int64(misc.NowMS()),
	}
	u := &User{
		PB:    pb,
		Items: make(map[int64]*def.Item),
		Mails: make(map[int64]bool),
	}

	return u, u.Save()
}

func LoadUser(id int64) (*User, error) {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("GET", fmt.Sprintf("user:%d:account", id))
	conn.Send("GET", fmt.Sprintf("user:%d:name", id))
	conn.Send("GET", fmt.Sprintf("user:%d:icon", id))
	conn.Send("GET", fmt.Sprintf("user:%d:lvl", id))
	conn.Send("GET", fmt.Sprintf("user:%d:vip", id))
	conn.Send("GET", fmt.Sprintf("user:%d:status", id))
	conn.Send("GET", fmt.Sprintf("user:%d:create_time", id))
	conn.Send("GET", fmt.Sprintf("user:%d:update_time", id))
	conn.Send("SMEMBERS", fmt.Sprintf("user:%d:items", id))
	conn.Send("SMEMBERS", fmt.Sprintf("user:%d:mails", id))
	res, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	account, err := redis.Int64(res[0], nil)
	if err != nil {
		return nil, err
	}

	name, err := redis.String(res[1], nil)
	if err != nil {
		return nil, err
	}

	icon, err := redis.Int(res[2], nil)
	if err != nil {
		return nil, err
	}

	lvl, err := redis.Int(res[3], nil)
	if err != nil {
		return nil, err
	}

	vip, err := redis.Int(res[4], nil)
	if err != nil {
		return nil, err
	}

	status, err := redis.Int(res[5], nil)
	if err != nil {
		return nil, err
	}

	createTime, err := redis.Int64(res[6], nil)
	if err != nil {
		return nil, err
	}

	updateTime, err := redis.Int64(res[7], nil)
	if err != nil {
		return nil, err
	}

	itemList, err := redis.Values(res[9], nil)
	if err != nil {
		return nil, err
	}

	mailList, err := redis.Values(res[9], nil)
	if err != nil {
		return nil, err
	}

	pb := &def.User{
		Id:         proto.Int64(id),
		Account:    proto.Int64(account),
		Name:       proto.String(name),
		Icon:       proto.Int32(int32(icon)),
		Lvl:        proto.Int32(int32(lvl)),
		Vip:        proto.Int32(int32(vip)),
		Status:     proto.Int32(int32(status)),
		CreateTime: proto.Int64(createTime),
		UpdateTime: proto.Int64(updateTime),
	}

	user := &User{
		PB:    pb,
		Items: make(map[int64]*def.Item),
		Mails: make(map[int64]bool),
	}

	for _, v := range itemList {
		id, _ := redis.Int64(v, nil)
		item, err := LoadItem(id)
		if err != nil {
			continue
		}
		user.Items[id] = item
	}

	for _, v := range mailList {
		id, err := redis.Int64(v, nil)
		if err != nil {
			continue
		}
		user.Mails[id] = true
	}

	return user, nil
}

func (u *User) ID() int64 {
	if u.PB != nil {
		return u.PB.GetId()
	}
	return 0
}

func (u *User) AccID() int64 {
	if u.PB != nil {
		return u.PB.GetAccount()
	}
	return 0
}

func (u *User) AddMail(id int64) error {
	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SADD", fmt.Sprintf("user:%d:mails", u.PB.Id), id); err != nil {
		return err
	}

	u.Mails[id] = true
	return nil
}

func (u *User) Save() error {
	conn := pool.Get()
	defer conn.Close()

	id := u.ID()
	pb := u.PB

	conn.Send("MULTI")

	conn.Send("SET", fmt.Sprintf("user:account:%d", pb.GetAccount()), id)
	conn.Send("SET", fmt.Sprintf("user:%d:account", id), pb.GetAccount())
	conn.Send("SET", fmt.Sprintf("user:%d:name", id), pb.GetName())
	conn.Send("SET", fmt.Sprintf("user:%d:icon", id), pb.GetIcon())
	conn.Send("SET", fmt.Sprintf("user:%d:lvl", id), pb.GetLvl())
	conn.Send("SET", fmt.Sprintf("user:%d:vip", id), pb.GetVip())
	conn.Send("SET", fmt.Sprintf("user:%d:status", id), pb.GetStatus())
	conn.Send("SET", fmt.Sprintf("user:%d:create_time", id), pb.GetCreateTime())
	conn.Send("SET", fmt.Sprintf("user:%d:update_time", id), pb.GetUpdateTime())
	for _, item := range u.Items {
		conn.Send("SADD", fmt.Sprintf("user:%d:items", id), item.Id)
	}
	for mailId, _ := range u.Mails {
		conn.Send("SADD", fmt.Sprintf("user:%d:mails", id), mailId)
	}
	_, err := redis.Values(conn.Do("EXEC"))
	return err
}
