package account

import (
	"fmt"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/labstack/echo"
	"gohive/internal/pb/def"
	"gohive/internal/pb/msg"
	"gohive/server/db"
)

func BindPhoneHandler(c echo.Context) error {
	req := &msg.BindPhone_Req{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	if db.FindAccountByPhone(req.GetPhone()) != nil {
		res := &msg.BindAccount_Res{
			Code: proto.Int32(int32(def.RC_ACCOUNT_PHONE_INVALID)),
		}
		return c.JSON(http.StatusOK, res)
	}

	fmt.Printf("bindPhone: %d %s %s\n", req.GetAccountId(), req.GetPhone(), req.GetCode())

	acc := db.FindAccountByID(req.GetAccountId())
	if acc == nil {
		res := &msg.BindAccount_Res{
			Code: proto.Int32(int32(def.RC_ACCOUNT_NO_FOUND)),
		}
		return c.JSON(http.StatusOK, res)
	}

	if !db.CheckAccount(acc, req.GetPassword()) {
		res := &msg.BindAccount_Res{
			Code: proto.Int32(int32(def.RC_AUTH_INVALID_USER)),
		}
		return c.JSON(http.StatusOK, res)
	}

	if acc.GetPhone() != "" {
		res := &msg.BindAccount_Res{
			Code: proto.Int32(int32(def.RC_ACCOUNT_NO_FOUND)),
		}
		return c.JSON(http.StatusOK, res)
	}

	ret := db.BindPhone(acc, req.GetPhone(), req.GetCode())
	res := &msg.BindAccount_Res{
		Code: proto.Int32(ret),
	}
	return c.JSON(http.StatusOK, res)
}
