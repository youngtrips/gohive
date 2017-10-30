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

func ResetPasswordHandler(c echo.Context) error {
	req := &msg.ResetPassword_Req{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	code := int32(def.RC_OK)

	if acc := db.FindAccountByID(req.GetAccountId()); acc == nil {
		code = int32(def.RC_ACCOUNT_NO_FOUND)
	} else {
		code = db.ResetAccountPassword(acc, req.GetNewPassword(), req.GetCode())
	}
	res := &msg.ResetPassword_Res{Code: proto.Int32(code)}
	return c.JSON(http.StatusOK, res)
}
