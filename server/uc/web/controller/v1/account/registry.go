package account

import (
	"fmt"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/labstack/echo"
	"gohive/internal/pb/def"
	"gohive/internal/pb/msg"
	"gohive/server/db"
	"gohive/server/uc/entity"
)

func RegistryHandler(c echo.Context) error {
	req := &msg.Registry_Req{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	res := &msg.Auth_Res{
		Code: proto.Int32(int32(def.RC_OK)),
	}

	id, _ := db.GetAccountId(req.GetUsername())
	if id > 0 {
		res.Code = proto.Int32(int32(def.RC_ACCOUNT_CREATE_DUPLICATE))
	} else {
		if acc, err := db.CreateAccount(req.GetUsername(), req.GetPassword()); err != nil {
			res.Code = proto.Int32(int32(def.RC_ACCOUNT_CREATE_FAILED))
		} else {
			service := entity.GetService()
			if service == nil {
				res.Code = proto.Int32(int32(def.RC_COMMON_SERVER_NOT_READY))
			} else {
				if token, err := db.GenToken(acc, service.Key); err != nil {
					fmt.Printf("gen token failed: %s\n", err)
					return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
				} else {
					res.Token = proto.String(token)
					res.Addr = proto.String(service.Addr)
				}
			}
		}
	}
	return c.JSON(http.StatusOK, res)
}
