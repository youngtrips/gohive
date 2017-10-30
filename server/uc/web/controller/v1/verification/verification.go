package verification

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/labstack/echo"
	"gohive/internal/pb/def"
	"gohive/internal/pb/msg"
	"gohive/sdk/aliyun/sms"
	"gohive/server/db"
)

func GetCodeHandler(c echo.Context) error {
	req := &msg.VerificationCode_Req{}
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%s", err))
	}

	fmt.Println(req)
	code := def.RC_OK

	verifyInfo, err := db.GenVerification(req.GetPhone())
	if err != nil {
		fmt.Println(err)
		code = def.RC_SERVER_INTERNAL_ERROR
	}
	fmt.Println(verifyInfo)
	log.Info("phone: ", req.GetPhone())
	if err := sms.SendCode(req.GetPhone(), verifyInfo.Code); err != nil {
		log.Error("sms.SendCode: ", err)
		code = def.RC_SERVER_INTERNAL_ERROR
	}
	res := &msg.VerificationCode_Res{
		Code: proto.Int32(int32(code)),
	}
	return c.JSON(http.StatusOK, res)
}
