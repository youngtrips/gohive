package db

import (
	"fmt"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"gohive/internal/misc"
	"gohive/internal/pb/def"
)

const (
	MAIL_ID_OFFSET = 10000
)

func GenMailID() (int64, error) {
	conn := pool.Get()
	defer conn.Close()

	id, err := redis.Int64(conn.Do("INCR", "mail:count"))
	if err != nil {
		return 0, err
	}
	return MAIL_ID_OFFSET + id, nil
}

func LoadMail(id int64) (*def.Mail, error) {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("GET", fmt.Sprintf("mail:%d:type", id))
	conn.Send("GET", fmt.Sprintf("mail:%d:from", id))
	conn.Send("GET", fmt.Sprintf("mail:%d:to", id))
	conn.Send("GET", fmt.Sprintf("mail:%d:status", id))
	conn.Send("GET", fmt.Sprintf("mail:%d:title", id))
	conn.Send("GET", fmt.Sprintf("mail:%d:content", id))
	conn.Send("GET", fmt.Sprintf("mail:%d:create_time", id))
	conn.Send("GET", fmt.Sprintf("mail:%d:update_time", id))
	conn.Send("GET", fmt.Sprintf("mail:%d:items", id))
	res, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	typ, err := redis.Int(res[0], nil)
	if err != nil {
		return nil, err
	}

	from, err := redis.Int64(res[1], nil)
	if err != nil {
		return nil, err
	}

	to, err := redis.Int64(res[2], nil)
	if err != nil {
		return nil, err
	}

	status, err := redis.Int(res[3], nil)
	if err != nil {
		return nil, err
	}

	title, err := redis.String(res[4], nil)
	if err != nil {
		return nil, err
	}

	content, err := redis.String(res[5], nil)
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

	items, err := redis.Int64Map(res[8], nil)
	if err != nil {
		return nil, err
	}

	mail := &def.Mail{
		Id:         proto.Int64(id),
		Type:       proto.Int32(int32(typ)),
		From:       proto.Int64(from),
		To:         proto.Int64(to),
		Title:      proto.String(title),
		Content:    proto.String(content),
		Status:     proto.Int32(int32(status)),
		CreateTime: proto.Int64(createTime),
		UpdateTime: proto.Int64(updateTime),
		Items:      make([]*def.Item, 0),
	}
	for idxS, num := range items {
		idx, err := strconv.Atoi(idxS)
		if err != nil {
			continue
		}
		item := &def.Item{
			Id:  proto.Int64(0),
			Idx: proto.Int32(int32(idx)),
			Num: proto.Int64(int64(num)),
		}

		mail.Items = append(mail.Items, item)
	}
	return mail, nil
}

func SendMail(from *User, to *User, typ, title string, content string, items []*def.Item) (*def.Mail, error) {
	mailId, err := GenMailID()
	if err != nil {
		return nil, err
	}
	mail := &def.Mail{
		Id:         proto.Int64(mailId),
		From:       proto.Int64(from.ID()),
		To:         proto.Int64(to.ID()),
		Status:     proto.Int32(int32(def.MAIL_STATUS_NEW)),
		Title:      proto.String(title),
		Content:    proto.String(content),
		CreateTime: proto.Int64(misc.NowMS()),
		UpdateTime: proto.Int64(misc.NowMS()),
		Items:      items,
	}
	if err := SaveMail(mail); err != nil {
		return nil, err
	}
	if err := to.AddMail(mailId); err != nil {
		DelMail(mailId)
		return nil, err
	}
	return mail, nil
}

func SaveMail(mail *def.Mail) error {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("GET", fmt.Sprintf("mail:%d:id", mail.GetId()), mail.GetId())
	conn.Send("GET", fmt.Sprintf("mail:%d:type", mail.GetId()), mail.GetType())
	conn.Send("GET", fmt.Sprintf("mail:%d:from", mail.GetId()), mail.GetFrom())
	conn.Send("GET", fmt.Sprintf("mail:%d:to", mail.GetId()), mail.GetTo())
	conn.Send("GET", fmt.Sprintf("mail:%d:status", mail.GetId()), mail.GetStatus())
	conn.Send("GET", fmt.Sprintf("mail:%d:title", mail.GetId()), mail.GetTitle())
	conn.Send("GET", fmt.Sprintf("mail:%d:content", mail.GetId()), mail.GetContent())
	conn.Send("GET", fmt.Sprintf("mail:%d:create_time", mail.GetId()), mail.GetCreateTime())
	conn.Send("GET", fmt.Sprintf("mail:%d:update_time", mail.GetId()), mail.GetUpdateTime())
	if mail.Items != nil {
		for _, item := range mail.Items {
			conn.Send("HSET", fmt.Sprintf("mail:%d:items", mail.GetId()), item.GetIdx(), item.GetNum())
		}
	}
	_, err := redis.Values(conn.Do("EXEC"))
	return err
}

func DelMail(id int64) error {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("DEL", fmt.Sprintf("mail:%d:id", id))
	conn.Send("DEL", fmt.Sprintf("mail:%d:type", id))
	conn.Send("DEL", fmt.Sprintf("mail:%d:from", id))
	conn.Send("DEL", fmt.Sprintf("mail:%d:to", id))
	conn.Send("DEL", fmt.Sprintf("mail:%d:status", id))
	conn.Send("DEL", fmt.Sprintf("mail:%d:title", id))
	conn.Send("DEL", fmt.Sprintf("mail:%d:content", id))
	conn.Send("DEL", fmt.Sprintf("mail:%d:create_time", id))
	conn.Send("DEL", fmt.Sprintf("mail:%d:update_time", id))
	conn.Send("DEL", fmt.Sprintf("mail:%d:items", id))
	_, err := redis.Values(conn.Do("EXEC"))
	return err
}
