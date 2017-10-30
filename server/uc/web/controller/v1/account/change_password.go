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

func ChangePasswordHandler(c echo.Context) error {
	req := &msg.ChangePassword_Req{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	code := int32(def.RC_OK)
	acc := db.FindAccountByID(req.GetAccountId())
	if acc == nil {
		code = int32(def.RC_ACCOUNT_NO_FOUND)
	} else {
		fmt.Println("old passworld: ", req.GetPassword())
		fmt.Println("new passworld: ", req.GetNewPassword())
		if !db.CheckAccount(acc, req.GetPassword()) {
			code = int32(def.RC_AUTH_INVALID_USER)
		} else {
			code = db.ChangeAccountPassword(acc, req.GetNewPassword())
		}
	}
	res := &msg.ChangePassword_Res{Code: proto.Int32(int32(code))}
	return c.JSON(http.StatusOK, res)
}
