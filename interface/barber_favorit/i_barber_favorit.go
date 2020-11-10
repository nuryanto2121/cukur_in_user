package ibarberfavorit

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Repository interface {
	GetDataBy(BarberId int, UserID int) (result *models.BarberFavorit, err error)
	GetList(queryparam models.ParamList) (result []*models.BarberFavoritList, err error)
	Create(data *models.BarberFavorit) (err error)
	Update(ID int, data interface{}) (err error)
	Delete(BarberId int, UserID int) (err error)
	Count(queryparam models.ParamList) (result int, err error)
}
type Usecase interface {
	// GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.AddBarberFavorit) (err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	// Update(ctx context.Context, Claims util.Claims, ID int, data models.UpdateUser) (err error)
	// Delete(ctx context.Context, ID int) (err error)
}
