package reponotification

import (
	"fmt"
	inotification "nuryanto2121/cukur_in_user/interface/notification"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoNotification struct {
	Conn *gorm.DB
}

func NewRepoNotification(Conn *gorm.DB) inotification.Repository {
	return &repoNotification{Conn}
}

func (db *repoNotification) GetDataBy(ID int) (result *models.Notification, err error) {
	var (
		logger        = logging.Logger{}
		mNotification = &models.Notification{}
	)
	query := db.Conn.Where("notification_id = ? ", ID).Find(mNotification)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mNotification, nil
}

func (db *repoNotification) GetList(UserID int, queryparam models.ParamListGeo) (result []*models.NotificationList, err error) {

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

	sSql := fmt.Sprintf(`
	select * from (
		select 	
			n.notification_id ,		n.notification_date ,
			n.notification_status, 	n.notification_type ,
			n.user_id, 				n.title ,
			n.descs , 				n.link_id,
			x.barber_id,			x.barber_name,
			x.order_id,				x.order_no,
			x.status,				x.order_date,
			x.price,				x.distance,
			x.barber_rating
		from notification n left join
		(SELECT oh.order_id ,oh.order_no ,oh.order_date ,a.barber_id,a.barber_name,a.distance,a.barber_rating,
		   (select sum(order_d.price ) from order_d where order_d.order_id = oh.order_id ) as price ,
		   oh.user_id ,oh.status 
		   FROM order_h oh join fbarber_beranda_user_s(%f,%f,%d) a 
		   on oh.barber_id = a.barber_id) x
	   on n.link_id = x.order_id) z
	
	`, queryparam.Latitude, queryparam.Longitude, UserID)

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and (lower(notification_status) LIKE ?)"
		} else {
			sWhere += "(lower(notification_status) LIKE ?)"
		}
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		// query = db.Conn.Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
		query = db.Conn.Raw(sSql, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		// query = db.Conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
		query = db.Conn.Raw(sSql).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	}

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

func (db *repoNotification) Create(data *models.Notification) error {
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
func (db *repoNotification) Update(ID int, data map[string]interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.Notification{}).Where("notification_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoNotification) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where("notification_id = ?", ID).Delete(&models.Notification{})
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoNotification) Count(UserID int, queryparam models.ParamListGeo) (result int, err error) {
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
			sWhere += " and (lower(notification_status) LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower(notification_status) LIKE ? )" //queryparam.Search
		}
		query = db.Conn.Model(&models.Notification{}).Where(sWhere, queryparam.Search).Count(&result)
	} else {
		query = db.Conn.Model(&models.Notification{}).Where(sWhere).Count(&result)
	}
	// end where

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
