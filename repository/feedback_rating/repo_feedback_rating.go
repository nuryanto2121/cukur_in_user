package repofeedbackrating

import (
	"fmt"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"

	ifeedbackrating "nuryanto2121/cukur_in_user/interface/feedback_rating"

	"github.com/jinzhu/gorm"
)

type repoFeedbackRating struct {
	Conn *gorm.DB
}

func NewRepoFeedbackRating(Conn *gorm.DB) ifeedbackrating.Repository {
	return &repoFeedbackRating{Conn}
}

func (db *repoFeedbackRating) GetDataBy(OrderID int) (result *models.FeedbackRating, err error) {
	var (
		logger = logging.Logger{}
		data   = &models.FeedbackRating{}
	)
	query := db.Conn.Raw(`select * from feedback_rating where order_id = ? `, OrderID).Scan(&data) //Find(&result)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return result, models.ErrNotFound
		}
		return result, err
	}
	return data, nil
}
func (db *repoFeedbackRating) GetList(queryparam models.ParamList) (result []*models.FeedbackRating, err error) {
	var (
		pageNum  = 0
		pageSize = setting.FileConfigSetting.App.PageSize
		sWhere   = ""
		logger   = logging.Logger{}
		orderBy  = queryparam.SortField
		query    *gorm.DB
	)
	// pagination
	if queryparam.Page > 0 {
		pageNum = (queryparam.Page - 1) * queryparam.PerPage
	}
	if queryparam.PerPage > 0 {
		pageSize = queryparam.PerPage
	}
	//end pagination

	// Order
	if queryparam.SortField != "" {
		orderBy = queryparam.SortField
	}
	//end Order by

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and (lower(comment) LIKE ? OR lower(descs) LIKE ?)" //+ queryparam.Search
		} else {
			sWhere += "(lower(comment) LIKE ? OR lower(descs) LIKE ?)" //queryparam.Search
		}
		query = db.Conn.Where(sWhere, queryparam.Search, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		query = db.Conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	}

	// end where

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (db *repoFeedbackRating) Create(data *models.FeedbackRating) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Create(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoFeedbackRating) Update(ID int, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.FeedbackRating{}).Where("id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoFeedbackRating) Delete(BarberId int, UserID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("order_id = ?", ID).Delete(&models.OrderH{})
	query := db.Conn.Exec("Delete From barber_favorit WHERE barber_id = ? and user_id = ?", BarberId, UserID)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoFeedbackRating) Count(queryparam models.ParamList) (result int, err error) {
	var (
		sWhere = ""
		logger = logging.Logger{}
		query  *gorm.DB
	)
	result = 0

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(b.barber_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(b.barber_name) LIKE ?" //queryparam.Search
		}
		query = db.Conn.Table("barber b ").Select(`
		b.barber_id,b.barber_cd,b.barber_name,
		b.address,b.latitude,b.longitude,
		b.operation_start,b.operation_end,
		b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type,true as is_favorit
		`).Joins(`
		join barber_favorit a
			on a.barber_id = b.barber_id 
		`).Joins(`
		left join sa_file_upload c
			on b.file_id = c.file_id 
	`).Where(sWhere, queryparam.Search).Count(&result)
	} else {
		query = db.Conn.Table("barber b ").Select(`
		b.barber_id,b.barber_cd,b.barber_name,
		b.address,b.latitude,b.longitude,
		b.operation_start,b.operation_end,
		b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type,true as is_favorit
		`).Joins(`
		join barber_favorit a
			on a.barber_id = b.barber_id 
		`).Joins(`
		left join sa_file_upload c
			on b.file_id = c.file_id 
	`).Where(sWhere).Count(&result)
	}
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
