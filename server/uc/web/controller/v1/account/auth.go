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

func AuthHandler(c echo.Context) error {
	req := &msg.Auth_Req{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	code := def.RC_OK
	token := ""
	addr := ""
	var acc *def.Account
	if req.GetUsername() == "" && req.GetDeviceId() == "" {
		code = def.RC_AUTH_INVALID_USER
	} else if req.GetUsername() != "" {
		fmt.Println("get account by name: ", req.GetUsername())
		acc = db.FindAccount(req.GetUsername())
		if acc == nil {
			code = def.RC_AUTH_INVALID_USER
		} else {
			if !db.CheckAccount(acc, req.GetPassword()) {
				code = def.RC_AUTH_INVALID_USER
			}
		}
	} else if req.GetDeviceId() != "" {
		fmt.Println("get account by deviceId: ", req.GetDeviceId())
		acc = db.GetAnonymousAccount(req.GetDeviceId())
	}
	if acc != nil {
		service := entity.GetService()
		if service == nil {
			code = def.RC_COMMON_SERVER_NOT_READY
		} else {
			var err error
			token, err = db.GenToken(acc, service.Key)
			if err != nil {
				fmt.Printf("gen token failed: %s\n", err)
				return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
			}
			addr = service.Addr
		}
	}

	res := &msg.Auth_Res{
		Code:  proto.Int32(int32(code)),
		Token: proto.String(token),
		Addr:  proto.String(addr),
	}
	return c.JSON(http.StatusOK, res)
}
