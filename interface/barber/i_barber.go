package ibarber

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Repository interface {
	GetDataBy(ID int) (result *models.Barber, err error)
	GetDataByList(ID int, UserID int, GeoBarber models.GeoBarber) (result *models.BarbersList, err error)
	GetDataFirst(OwnerID int, BarberID int) (result *models.Barber, err error)
	GetList(UserID int, queryparam models.ParamListGeo) (result []*models.BarbersList, err error)
	GetScheduleTime(BarberID int) (result interface{}, err error)
	Create(data *models.Barber) (err error)
	Update(ID int, data interface{}) (err error)
	Delete(ID int) (err error)
	Count(UserID int, queryparam models.ParamListGeo) (result int, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int, GeoBarber models.GeoBarber) (result interface{}, err error)
	GetDataFirst(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.BarbersPost) error
	Update(ctx context.Context, Claims util.Claims, ID int, data models.BarbersPost) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID int) (err error)
}
