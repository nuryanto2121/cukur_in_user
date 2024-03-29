package iauth

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Repository interface {
	GetDataLogin(ctx context.Context, Account string) (models.DataLogin, error)
	ChangePassword(ctx context.Context, data interface{}) (err error)
	Register(ctx context.Context, dataUser models.SsUser) error
}
type Usecase interface {
	Logout(ctx context.Context, Claims util.Claims, Token string) (err error)
	Login(ctx context.Context, dataLogin *models.LoginForm) (output interface{}, err error)
	ForgotPassword(ctx context.Context, dataForgot *models.ForgotForm) (result string, err error)
	GenOTP(ctx context.Context, dataForgot *models.ForgotForm) (result interface{}, err error)
	ResetPassword(ctx context.Context, dataReset *models.ResetPasswd) (err error)
	Register(ctx context.Context, dataRegister models.RegisterForm) (output interface{}, err error)
	Verify(ctx context.Context, dataVeriry models.VerifyForm) (err error)
	VerifyRegister(ctx context.Context, dataVerify models.VerifyForm) (output interface{}, err error)
}
