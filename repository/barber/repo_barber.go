package repobarber

import (
	"fmt"
	ibarber "nuryanto2121/cukur_in_user/interface/barber"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoBarber struct {
	Conn *gorm.DB
}

func NewRepoBarber(Conn *gorm.DB) ibarber.Repository {
	return &repoBarber{Conn}
}

func (db *repoBarber) GetDataBy(ID int, GeoBarber models.GeoBarber) (result *models.BarbersList, err error) {
	var (
		logger  = logging.Logger{}
		mBarber = &models.BarbersList{}
	)
	sSql := fmt.Sprintf(`
		SELECT * FROM fbarber_beranda_user_s(%f,%f)
		WHERE barber_id = ?
	`, GeoBarber.Latitude, GeoBarber.Longitude)

	// query := db.Conn.Where("barber_id = ? ", ID).Find(mBarber)
	query := db.Conn.Raw(sSql, ID).First(&mBarber)

	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mBarber, nil
}
func (db *repoBarber) GetDataFirst(OwnerID int, BarberID int) (result *models.Barber, err error) {
	var (
		logger  = logging.Logger{}
		mBarber = &models.Barber{}
		sQuery  = ""
	)
	// query := db.Conn.First(&mBarber)
	if BarberID == 0 {
		sQuery = `SELECT * FROM barber where is_active = true and owner_id = ? 
		order by barber_id 
		limit 1`

		query := db.Conn.Raw(sQuery, OwnerID).Scan(&mBarber)
		logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
		err = query.Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, models.ErrNotFound
			}
			return nil, err
		}
	} else {
		sQuery = `SELECT * FROM barber where is_active = true and owner_id = ? AND barber_id = ?
		limit 1`

		query := db.Conn.Raw(sQuery, OwnerID, BarberID).Scan(&mBarber)
		logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
		err = query.Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, models.ErrNotFound
			}
			return nil, err
		}
	}

	return mBarber, nil
}
func (db *repoBarber) GetList(queryparam models.ParamListGeo) (result []*models.BarbersList, err error) {

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
		SELECT * FROM fbarber_beranda_user_s(%f,%f)
	`, queryparam.Latitude, queryparam.Longitude)

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(barber_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(barber_name) LIKE ?" //queryparam.Search
		}
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		query = db.Conn.Raw(sSql, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

	} else {
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
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
func (db *repoBarber) Create(data *models.Barber) error {
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
func (db *repoBarber) Update(ID int, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.Barber{}).Where("barber_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarber) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("barber_id = ?", ID).Delete(&models.Barber{})
	query := db.Conn.Exec("Delete From barber_collection WHERE barber_id = ?", ID)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarber) Count(queryparam models.ParamListGeo) (result int, err error) {

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
	sSql := fmt.Sprintf(`
	SELECT count(*) as cnt FROM fbarber_beranda_user_s(%f,%f)
`, queryparam.Latitude, queryparam.Longitude)
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
