package inotification

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Repository interface {
	GetDataBy(ID int) (result *models.Notification, err error)
	GetList(UserID int, queryparam models.ParamListGeo) (result []*models.NotificationList, err error)
	Create(data *models.Notification) (err error)
	Update(ID int, data map[string]interface{}) (err error)
	Delete(ID int) (err error)
	Count(UserID int, queryparam models.ParamListGeo) (result int, err error)
}

type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result *models.Notification, err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamListGeo) (result models.ResponseModelList, err error)
	Create(ctx context.Context, Claims util.Claims, TokenFCM string, data *models.AddNotification) (err error)
	Update(ctx context.Context, Claims util.Claims, ID int, data *models.StatusNotification) (err error)
	Delete(ctx context.Context, Claims util.Claims, ID int) (err error)
	GetCountNotif(ctx context.Context, Claims util.Claims) (result interface{}, err error)
}
