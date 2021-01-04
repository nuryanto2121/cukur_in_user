package icorder

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Repository interface {
	GetDataBy(ID int, GeoUser models.GeoBarber) (result models.OrderDGet, err error)
	GetList(UserID int, queryparam models.ParamListGeo) (result []*models.OrderList, err error)
	Create(data *models.OrderH) (err error)
	Update(ID int, data interface{}) (err error)
	Delete(ID int) (err error)
	Count(UserID int, queryparam models.ParamListGeo) (result int, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int, GeoUser models.GeoBarber) (result interface{}, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.OrderPost) error
	Update(ctx context.Context, Claims util.Claims, ID int, data models.OrderStatus) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID int) (err error)
}
