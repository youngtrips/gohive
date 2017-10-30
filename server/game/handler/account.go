package handler

import (
	"gohive/internal/pb/msg"
)

func OnAuthReq(req *msg.Auth_Req) *msg.Auth_Res {
	res := &msg.Auth_Res{}
	return res
}

func OnAuthRes(res *msg.Auth_Res) {
}

func OnRegistryReq(req *msg.Registry_Req) *msg.Registry_Res {
	res := &msg.Registry_Res{}
	return res
}

func OnRegistryRes(res *msg.Registry_Res) {
}

func OnBindAccountReq(req *msg.BindAccount_Req) *msg.BindAccount_Res {
	res := &msg.BindAccount_Res{}
	return res
}

func OnBindAccountRes(res *msg.BindAccount_Res) {
}

func OnBindPhoneReq(req *msg.BindPhone_Req) *msg.BindPhone_Res {
	res := &msg.BindPhone_Res{}
	return res
}

func OnBindPhoneRes(res *msg.BindPhone_Res) {
}

func OnVerificationCodeReq(req *msg.VerificationCode_Req) *msg.VerificationCode_Res {
	res := &msg.VerificationCode_Res{}
	return res
}

func OnVerificationCodeRes(res *msg.VerificationCode_Res) {
}

func OnResetPasswordReq(req *msg.ResetPassword_Req) *msg.ResetPassword_Res {
	res := &msg.ResetPassword_Res{}
	return res
}

func OnResetPasswordRes(res *msg.ResetPassword_Res) {
}

func OnChangePasswordReq(req *msg.ChangePassword_Req) *msg.ChangePassword_Res {
	res := &msg.ChangePassword_Res{}
	return res
}

func OnChangePasswordRes(res *msg.ChangePassword_Res) {
}
