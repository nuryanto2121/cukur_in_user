package icapster

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Repository interface {
	GetDataBy(ID int, GeoBarber models.GeoBarber) (result *models.CapsterList, err error)
	GetListFileCapter(ID int) (result []*models.SaFileOutput, err error)
	GetList(queryparam models.ParamListGeo) (result []*models.CapsterList, err error)
	Count(queryparam models.ParamListGeo) (result int, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int, GeoBarber models.GeoBarber) (result interface{}, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error)
}
