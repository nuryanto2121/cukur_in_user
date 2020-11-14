package repobarbercapster

import (
	"fmt"
	ibarbercapster "nuryanto2121/cukur_in_user/interface/barber_capster"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoBarberCapster struct {
	Conn *gorm.DB
}

func NewRepoBarberCapster(Conn *gorm.DB) ibarbercapster.Repository {
	return &repoBarberCapster{Conn}
}

func (db *repoBarberCapster) GetDataBy(ID int) (result *models.BarberCapster, err error) {
	var (
		logger         = logging.Logger{}
		mBarberCapster = &models.BarberCapster{}
	)
	query := db.Conn.Where("barber_id = ? ", ID).Find(mBarberCapster)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mBarberCapster, nil
}
func (db *repoBarberCapster) GetList(queryparam models.ParamListGeo) (result []*models.CapsterList, err error) {

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
		SELECT * FROM fbarber_capster_s(%f,%f)
	`, queryparam.Latitude, queryparam.Longitude)

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(capster_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(capster_name) LIKE ?" //queryparam.Search
		}
		query = db.Conn.Raw(sSql).Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

	} else {
		query = db.Conn.Raw(sSql).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

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
func (db *repoBarberCapster) Create(data *models.BarberCapster) error {
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
func (db *repoBarberCapster) Update(ID int, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.BarberCapster{}).Where("barber_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberCapster) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("barber_id = ?", ID).Delete(&models.BarberCapster{})
	query := db.Conn.Exec("Delete From barber_capster WHERE barber_id = ?", ID)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberCapster) DeleteByCapster(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("barber_id = ?", ID).Delete(&models.BarberCapster{})
	query := db.Conn.Exec("Delete From barber_capster WHERE capster_id = ?", ID)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberCapster) Count(queryparam models.ParamListGeo) (result int, err error) {
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
	sSql := fmt.Sprintf(`
	SELECT * FROM fbarber_capster_s(%f,%f)
`, queryparam.Latitude, queryparam.Longitude)
	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(capster_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(capster_name) LIKE ?" //queryparam.Search
		}
		query = db.Conn.Raw(sSql).Where(sWhere, queryparam.Search).Count(&result)

	} else {
		query = db.Conn.Raw(sSql).Where(sWhere).Count(&result)
	}
	// end where

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
