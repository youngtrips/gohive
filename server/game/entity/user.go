package entity

import (
	log "github.com/Sirupsen/logrus"
	"gohive/server/db"
)

var (
	_users       map[int64]*db.User
	_account_ids map[int64]int64
)

func init() {
	_users = make(map[int64]*db.User)
	_account_ids = make(map[int64]int64)
}

func AddUser(user *db.User) {
	_users[user.ID()] = user
	_account_ids[user.AccID()] = user.ID()
}

func GetUser(id int64) *db.User {
	return _users[id]
}

func GetUserByAccount(accId int64) *db.User {
	userId, present := _account_ids[accId]
	if !present {
		return nil
	}
	return GetUser(userId)
}

func LoadUser(accId int64) *db.User {
	log.Infof("load user: accId[%d]", accId)
	user := GetUserByAccount(accId)
	if user != nil {
		return user
	}

	userId, err := db.GetUserId(accId)
	if err != nil {
		log.Info("getUserId by account failed: ", err)
		return nil
	}

	if user, err = db.LoadUser(userId); err != nil {
		log.Warn("loadUser failed: ", err)
	}
	if user == nil {
		user, err = db.CreateUser(accId)
		if err != nil {
			log.Error("create user failed: ", err)
			return nil
		}
	}
	AddUser(user)
	return user
}

func DelUser(accId int64) {
	user := GetUserByAccount(accId)
	if user != nil {
		delete(_users, user.ID())
		delete(_account_ids, user.AccID())
	}
}
