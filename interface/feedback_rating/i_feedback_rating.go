package ifeedbackrating

import (
	"context"
	"nuryanto2121/cukur_in_user/models"
	util "nuryanto2121/cukur_in_user/pkg/utils"
)

type Repository interface {
	GetDataBy(OrderID int) (result *models.FeedbackRating, err error)
	GetList(queryparam models.ParamList) (result []*models.FeedbackRating, err error)
	Create(data *models.FeedbackRating) (err error)
	Update(ID int, data interface{}) (err error)
	Delete(BarberId int, UserID int) (err error)
	Count(queryparam models.ParamList) (result int, err error)
}
type Usecase interface {
	GetDataBy(ctx context.Context, Claims util.Claims, ID int) (result interface{}, err error)
	Create(ctx context.Context, Claims util.Claims, data *models.AddFeedbackRating) (err error)
	GetList(ctx context.Context, Claims util.Claims, queryparam models.ParamList) (result models.ResponseModelList, err error)
	Update(ctx context.Context, Claims util.Claims, ID int, data models.AddFeedbackRating) (err error)
	// Delete(ctx context.Context, ID int) (err error)
}
