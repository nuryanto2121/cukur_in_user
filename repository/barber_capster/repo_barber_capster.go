package repobarbercapster

import (
	"fmt"
	ibarbercapster "nuryanto2121/cukur_in_user/interface/barber_capster"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"

	"gorm.io/gorm"
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
	logger.Query(fmt.Sprintf("%v", query))
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
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		if orderBy != "" {
			sSql += fmt.Sprintf("\n order by %s", orderBy)
		}
		sSql += fmt.Sprintf("\n offset %d limit %d", pageNum, pageSize)
		query = db.Conn.Raw(sSql, queryparam.Search).Find(&result)
		// sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		// query = db.Conn.Raw(sSql, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

	} else {
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		if orderBy != "" {
			sSql += fmt.Sprintf("\n order by %s", orderBy)
		}
		sSql += fmt.Sprintf("\n offset %d limit %d", pageNum, pageSize)
		query = db.Conn.Raw(sSql).Find(&result)
		// sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		// query = db.Conn.Raw(sSql).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

	}

	logger.Query(fmt.Sprintf("%v", sSql)) //cath to log query string
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
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
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
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
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
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
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
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberCapster) Count(queryparam models.ParamListGeo) (result int, err error) {

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
	SELECT count(*) as cnt FROM fbarber_capster_s(%f,%f)
`, queryparam.Latitude, queryparam.Longitude)
	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(capster_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(capster_name) LIKE ?" //queryparam.Search
		}
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		query = db.Conn.Raw(sSql, queryparam.Search).First(&op)

	} else {
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		query = db.Conn.Raw(sSql).First(&op)
	}
	// end where

	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return op.Cnt, nil
}
