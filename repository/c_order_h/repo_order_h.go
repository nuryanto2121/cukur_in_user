package repoorderh

import (
	"fmt"
	iorder_h "nuryanto2121/cukur_in_user/interface/c_order_h"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoOrderH struct {
	Conn *gorm.DB
}

func NewRepoOrderH(Conn *gorm.DB) iorder_h.Repository {
	return &repoOrderH{Conn}
}

func (db *repoOrderH) GetDataBy(ID int, GeoUser models.GeoBarber) (result models.OrderDGet, err error) {
	var (
		logger = logging.Logger{}
		data   models.OrderDGet
	)
	sSql := fmt.Sprintf(`	
			select oh.order_id ,		oh.order_no, 
			oh.user_id ,		oh.status ,
			oh.order_date,	 	a.barber_id,
			a.barber_name,		a.capster_rating,
			a.distance,			a.barber_rating,
			a.capster_id,		a.capster_name,
			a.file_id,			a.file_name,
			a.file_path,		a.file_type,
			oh.from_apps,		(
				select sum(od.price) from order_d od 
				where od.order_id =oh.order_id 
			) as total_price
			from fbarber_capster_s(%f,%f) a
			join order_h oh on oh.barber_id = a.barber_id
			and oh.capster_id = a.capster_id
			AND oh.order_id = ?

	`, GeoUser.Latitude, GeoUser.Longitude)
	query := db.Conn.Raw(sSql, ID).Scan(&data) //Find(&result)
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
func (db *repoOrderH) GetList(UserID int, queryparam models.ParamListGeo) (result []*models.OrderList, err error) {

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

	sSql := fmt.Sprintf(`SELECT * FROM (
		SELECT oh.order_id ,oh.order_no ,oh.order_date ,a.barber_id,a.barber_name,a.distance,a.barber_rating,
		(select sum(order_d.price ) from order_d where order_d.order_id = oh.order_id ) as price ,
		oh.user_id ,oh.status 
		FROM order_h oh join fbarber_beranda_user_s(%f,%f,%d) a 
		on oh.barber_id = a.barber_id
		) xx`, queryparam.Latitude, queryparam.Longitude, UserID)

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(barber_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(barber_name) LIKE ?" //queryparam.Search
		}
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		fmt.Println(sSql)
		query = db.Conn.Raw(sSql, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		fmt.Println(sSql)
		query = db.Conn.Raw(sSql).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	}

	// end where

	// query := db.Conn.Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

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
func (db *repoOrderH) Create(data *models.OrderH) error {
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
func (db *repoOrderH) Update(ID int, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.OrderH{}).Where("order_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoOrderH) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("order_id = ?", ID).Delete(&models.OrderH{})
	query := db.Conn.Exec("Delete From order_h WHERE order_id = ?", ID)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoOrderH) Count(UserID int, queryparam models.ParamListGeo) (result int, err error) {
	type Results struct {
		Cnt int `json:"cnt"`
	}
	var (
		sWhere = ""
		logger = logging.Logger{}
		op     = &Results{}
		query  *gorm.DB
	)
	result = 0

	// WHERE
	if queryparam.InitSearch != "" {
		sWhere = queryparam.InitSearch
	}
	sSql := fmt.Sprintf(`SELECT count(*) as cnt FROM (
		SELECT oh.order_id ,oh.order_no ,oh.order_date ,a.barber_id,a.barber_name,a.distance,a.barber_rating,
		(select sum(order_d.price ) from order_d where order_d.order_id = oh.order_id ) as price ,
		oh.user_id ,oh.status 
		FROM order_h oh join fbarber_beranda_user_s(%f,%f,%d) a 
		on oh.barber_id = a.barber_id
		) xx`, queryparam.Latitude, queryparam.Longitude, UserID)

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(barber_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(barber_name) LIKE ?" //queryparam.Search
		}
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		query = db.Conn.Raw(sSql, queryparam.Search).First(&op)
	} else {
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		query = db.Conn.Raw(sSql).First(&op)
	}
	// end where

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return op.Cnt, nil
}
