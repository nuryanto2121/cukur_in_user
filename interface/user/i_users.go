package iusers

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Repository interface {
	GetDataBy(ID int) (result *models.SsUser, err error)
	GetByAccount(Account string) (result models.SsUser, err error)
	GetList(queryparam models.ParamList) (result []*models.SsUser, err error)
	Create(data *models.SsUser) (err error)
	Update(ID int, data interface{}) (err error)
	Delete(ID int) (err error)
	UpdatePasswordByEmail(Email string, Password string) error
	Count(queryparam models.ParamList) (result int, err error)
}
type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error)
	ChangePassword(ctx context.Context, Claims util.Claims, DataChangePwd models.ChangePassword) (err error)
	GetByEmailSaUser(ctx context.Context, email string) (result models.SsUser, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.SsUser) (err error)
	Update(ctx context.Context, Claims util.Claims, ID int, data models.UpdateUser) (err error)
	Delete(ctx context.Context, ID int) (err error)
}
