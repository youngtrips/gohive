package web

import (
	"gohive/server/uc/web/controller/v1/account"
	"gohive/server/uc/web/controller/v1/verification"

	"github.com/labstack/echo"
)

func route_init(e *echo.Echo) {
	g := e.Group("/v1")

	g.GET("/account/auth", account.AuthHandler)
	g.GET("/account/registry", account.RegistryHandler)
	g.GET("/account/bind/account", account.BindAccountHandler)
	g.GET("/account/bind/phone", account.BindPhoneHandler)
	g.GET("/account/password/reset", account.ResetPasswordHandler)
	g.GET("/account/password/change", account.ChangePasswordHandler)
	g.GET("/verification", verification.GetCodeHandler)

	g.POST("/account/auth", account.AuthHandler)
	g.POST("/account/registry", account.RegistryHandler)
	g.POST("/account/bind/account", account.BindAccountHandler)
	g.POST("/account/bind/phone", account.BindPhoneHandler)
	g.POST("/account/password/reset", account.ResetPasswordHandler)
	g.POST("/account/password/change", account.ChangePasswordHandler)
	g.POST("/verification", verification.GetCodeHandler)
}
