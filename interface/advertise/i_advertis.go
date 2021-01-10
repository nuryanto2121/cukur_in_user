package iadvertise

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Repository interface {
	GetDataBy(ID int) (result *models.Advertise, err error)
	GetList(queryparam models.ParamList) (result []*models.ListAdvertise, err error)
	Create(data *models.Advertise) (err error)
	Update(ID int, data map[string]interface{}) (err error)
	Delete(ID int) (err error)
	Count(queryparam models.ParamList) (result int, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.AddAdvertise) (err error)
	Update(ctx context.Context, Claims util.Claims, ID int, data *models.AddAdvertise) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID int) (err error)
}
