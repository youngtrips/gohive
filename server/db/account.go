package db

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"gohive/internal/misc"
	"gohive/internal/pb/def"
)

const (
	ACCOUNT_ID_OFFSET = 10000
	TOKEN_MAXN_TTL    = 3600 // 3600 seconds
)

var (
	r *rand.Rand
)

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GetAccountId(username string) (int64, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("GET", fmt.Sprintf("account:username:%s", username)))
}

func FindAccountByID(id int64) *def.Account {
	conn := pool.Get()
	defer conn.Close()

	/*
		       message Account {
		    required int64 id                = 1;
		    required string username         = 2;
		    required string password         = 3;
		    optional string device_id        = 4;
		    optional string phone            = 6;
		    optional int32 status            = 8;
		    optional int64 createtime        = 10;
		}
	*/

	conn.Send("MULTI")
	conn.Send("GET", fmt.Sprintf("account:%d:username", id))
	conn.Send("GET", fmt.Sprintf("account:%d:password", id))
	conn.Send("GET", fmt.Sprintf("account:%d:device_id", id))
	conn.Send("GET", fmt.Sprintf("account:%d:phone", id))
	r, err := redis.Strings(conn.Do("EXEC"))
	if err != nil {
		return nil
	}

	return &def.Account{
		Id:       proto.Int64(id),
		Username: proto.String(r[0]),
		Password: proto.String(r[1]),
		DeviceId: proto.String(r[2]),
		Phone:    proto.String(r[3]),
	}
}

func FindAccount(username string) *def.Account {
	id, err := GetAccountId(username)
	if err != nil {
		return nil
	}
	return FindAccountByID(id)
}

func SaveAccount(acc *def.Account) error {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SADD", "account:idlist", acc.GetId())
	conn.Send("SET", fmt.Sprintf("account:username:%s", acc.GetUsername()), acc.GetId())
	conn.Send("SET", fmt.Sprintf("account:%d:username", acc.GetId()), acc.GetUsername())
	conn.Send("SET", fmt.Sprintf("account:%d:password", acc.GetId()), acc.GetPassword())
	conn.Send("SET", fmt.Sprintf("account:%d:device_id", acc.GetId()), acc.GetDeviceId())
	conn.Send("SET", fmt.Sprintf("account:%d:phone", acc.GetId()), acc.GetPhone())
	_, err := conn.Do("EXEC")
	if err != nil {
		log.Error("SaveAccount: ", err)
	}
	return err
}

func GenAccountID() (int64, error) {
	conn := pool.Get()
	defer conn.Close()

	id, err := redis.Int64(conn.Do("INCR", "account:count"))
	if err != nil {
		return 0, err
	}
	return ACCOUNT_ID_OFFSET + id, nil
}

func genSalt() string {
	src := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	dst := make([]byte, 6)
	for i := 0; i < 6; i++ {
		j := r.Int31n(int32(len(src)))
		dst[i] = src[j]
	}
	return string(dst)
}

func genCode() string {
	src := []byte("0123456789")
	dst := make([]byte, 6)
	for i := 0; i < 6; i++ {
		j := r.Int31n(int32(len(src)))
		dst[i] = src[j]
	}
	return string(dst)
}

func CreateAccount(username string, password string) (*def.Account, error) {
	id, err := GenAccountID()
	if err != nil {
		return nil, err
	}

	salt := genSalt()

	log.Infof("createAccount: id=%d", id)

	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	password = fmt.Sprintf("%x", md5.Sum([]byte(password+salt)))
	password = base64.StdEncoding.EncodeToString([]byte(password + ";" + salt))

	log.Info("salt: ", salt)
	log.Info("pass: ", password)

	acc := &def.Account{
		Id:       proto.Int64(id),
		Username: proto.String(username),
		Password: proto.String(password),
	}
	return acc, SaveAccount(acc)
}

func getSalt(password string) (string, error) {
	ctx, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return "", err
	}
	strs := strings.Split(string(ctx), ";")
	if len(strs) != 2 {
		return "", nil
	}
	return strs[1], nil
}

func CheckAccount(acc *def.Account, password string) bool {
	salt, err := getSalt(acc.GetPassword())
	if err != nil {
		return false
	}

	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	password = fmt.Sprintf("%x", md5.Sum([]byte(password+salt)))
	password = base64.StdEncoding.EncodeToString([]byte(password + ";" + salt))

	log.Info("src password: ", acc.GetPassword())
	log.Info("dst password: ", password)

	return acc.GetPassword() == password
}

type MyClaims struct {
	AccountId int64 `json:"account_id"`
	jwt.StandardClaims
}

func GenToken(acc *def.Account, privateKey string) (string, error) {
	//https://dinosaurscode.xyz/go/2016/06/17/golang-jwt-authentication/
	expireToken := time.Now().Add(time.Hour * 2).Unix()

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", err
	}
	claims := MyClaims{
		acc.GetId(),
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "uc",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	log.Info("token: ", token)
	return token.SignedString(key)
}

func CheckToken(tokenStr string, pubKey string) (int64, int32) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwt.ParseRSAPublicKeyFromPEM([]byte(pubKey))
	})
	if err != nil {
		log.Error("CheckToken: ", err)
		return 0, int32(def.RC_LOGIN_TOKEN_INVALID)
	}

	if claims, ok := token.Claims.(MyClaims); ok && token.Valid {
		return claims.AccountId, int32(def.RC_OK)
	}

	return 0, int32(def.RC_LOGIN_TOKEN_INVALID)
}

func BindAccount(acc *def.Account, username string, password string) int32 {
	salt := genSalt()

	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	password = fmt.Sprintf("%x", md5.Sum([]byte(password+salt)))
	password = base64.StdEncoding.EncodeToString([]byte(password + ";" + salt))

	log.Info("salt: ", salt)
	log.Info("pass: ", password)

	acc.Username = proto.String(username)
	acc.Password = proto.String(password)
	log.Info(acc)
	if err := SaveAccount(acc); err != nil {
		log.Error(err)
		return int32(def.RC_SERVER_INTERNAL_ERROR)
	}
	return int32(def.RC_OK)
}

//	db.BindPhone(acc, req.GetPhone(), req.GetCode())
// verification:[phone]:code
// verification:[phone]:expired

type Verification struct {
	Code    string
	Expired int64
}

func GenVerification(phone string) (*Verification, error) {
	conn := pool.Get()
	defer conn.Close()

	code := genCode()
	expired := misc.NowMS() + 300000 // ttl = 5 min

	conn.Send("MULTI")
	conn.Send("SET", fmt.Sprintf("verification:%s:code", phone), code)
	conn.Send("SET", fmt.Sprintf("verification:%s:expired", phone), expired)
	_, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return nil, err
	}
	return &Verification{code, expired}, nil
}

func GetVerification(phone string) (*Verification, error) {
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("GET", fmt.Sprintf("verification:%s:code", phone))
	conn.Send("GET", fmt.Sprintf("verification:%s:expired", phone))
	r, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return nil, err
	}

	code, err := redis.String(r[0], nil)
	if err != nil {
		return nil, err
	}
	expired, err := redis.Int64(r[1], nil)
	if err != nil {
		return nil, err
	}
	return &Verification{code, expired}, nil
}

func BindPhone(acc *def.Account, phone string, code string) int32 {
	now := misc.NowMS()
	verifyInfo, err := GetVerification(phone)
	if err != nil {
		log.Info(err)
		return int32(def.RC_ACCOUNT_PHONE_INVALID)
	}
	if verifyInfo.Expired < now {
		return int32(def.RC_ACCOUNT_VERIFICATION_EXPIRED)
	}
	if verifyInfo.Code != code {
		return int32(def.RC_ACCOUNT_VERIFICATION_INVALID)
	}

	log.Info("bind: ", phone, " ", code)

	acc.Phone = proto.String(phone)
	if err := SaveAccount(acc); err != nil {
		log.Info(err)
		return int32(def.RC_SERVER_INTERNAL_ERROR)
	}
	return int32(def.RC_OK)
}

func FindAccountByPhone(phone string) *def.Account {
	conn := pool.Get()
	defer conn.Close()

	accId, err := redis.Int64(conn.Do("GET", fmt.Sprintf("account:phone:%s", phone)))
	if err != nil {
		return nil
	}
	return FindAccountByID(accId)
}

func ChangeAccountPassword(acc *def.Account, password string) int32 {
	salt := genSalt()

	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	password = fmt.Sprintf("%x", md5.Sum([]byte(password+salt)))
	password = base64.StdEncoding.EncodeToString([]byte(password + ";" + salt))

	log.Info("salt: ", salt)
	log.Info("pass: ", password)

	acc.Password = proto.String(password)
	if err := SaveAccount(acc); err != nil {
		log.Error(err)
		return int32(def.RC_SERVER_INTERNAL_ERROR)
	}
	return int32(def.RC_OK)
}

func ResetAccountPassword(acc *def.Account, password string, code string) int32 {
	log.Info("reset: ", acc, " ", password, " ", code)
	now := misc.NowMS()
	verifyInfo, err := GetVerification(acc.GetPhone())
	if err != nil {
		log.Info(err)
		return int32(def.RC_ACCOUNT_PHONE_INVALID)
	}
	log.Info(verifyInfo)
	if verifyInfo.Expired < now {
		return int32(def.RC_ACCOUNT_VERIFICATION_EXPIRED)
	}
	if verifyInfo.Code != code {
		return int32(def.RC_ACCOUNT_VERIFICATION_INVALID)
	}
	return ChangeAccountPassword(acc, password)
}

func getAnonymousAccountId(deviceId string) (int64, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("GET", fmt.Sprintf("anonymous:%s:account_id", deviceId)))
}

func saveAnonymousAccountId(deviceId string, id int64) {
	conn := pool.Get()
	defer conn.Close()
	conn.Do("SET", fmt.Sprintf("anonymous:%s:account_id", deviceId), id)
}

func GetAnonymousAccount(deviceId string) *def.Account {
	id, _ := getAnonymousAccountId(deviceId)
	var acc *def.Account
	if id > 0 {
		acc = FindAccountByID(id)
	}
	if acc == nil {
		id, err := GenAccountID()
		if err != nil {
			log.Error(err)
			return nil
		}
		acc = &def.Account{
			Id:       proto.Int64(id),
			Username: proto.String(""),
			Password: proto.String(""),
			DeviceId: proto.String(deviceId),
			Phone:    proto.String(""),
		}

		if err := SaveAccount(acc); err != nil {
			log.Error(err)
			return nil
		}
	}
	log.Info("acc: ", acc)
	saveAnonymousAccountId(deviceId, acc.GetId())
	return acc
}

/*
account:count id
account:userlist set(id)
account:email:[email] id

account:[id]:version number
account:[id]:email string
account:[id]:password string // md5(password..salt)
account:[id]:nickname string
account:[id]:lastlogin hashes
ip string
time string
account:[id]:history list(string)
account:[id]:available enum(open/locked/delete)


// Use the Send and Do methods to implement pipelined transactions.
//
//  c.Send("MULTI")
//  c.Send("INCR", "foo")
//  c.Send("INCR", "bar")
//  r, err := c.Do("EXEC")
//  fmt.Println(r) // prints [0, 1]
//
/

*/
