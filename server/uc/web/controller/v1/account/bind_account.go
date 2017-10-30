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

func BindAccountHandler(c echo.Context) error {
	req := &msg.BindAccount_Req{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	if db.FindAccount(req.GetUsername()) != nil {
		res := &msg.BindAccount_Res{
			Code: proto.Int32(int32(def.RC_ACCOUNT_CREATE_DUPLICATE)),
		}
		return c.JSON(http.StatusOK, res)
	}

	acc := db.FindAccountByID(req.GetAccountId())
	if acc == nil {
		res := &msg.BindAccount_Res{
			Code: proto.Int32(int32(def.RC_ACCOUNT_NO_FOUND)),
		}
		return c.JSON(http.StatusOK, res)
	}

	if acc.GetUsername() != "" {
		res := &msg.BindAccount_Res{
			Code: proto.Int32(int32(def.RC_ACCOUNT_NO_FOUND)),
		}
		return c.JSON(http.StatusOK, res)
	}

	ret := db.BindAccount(acc, req.GetUsername(), req.GetPassword())
	res := &msg.BindAccount_Res{
		Code: proto.Int32(ret),
	}
	return c.JSON(http.StatusOK, res)
}
