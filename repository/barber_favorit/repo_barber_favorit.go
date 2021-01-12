package repobarberfavorit

import (
	"fmt"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"

	ibarberfavorit "nuryanto2121/cukur_in_user/interface/barber_favorit"

	"gorm.io/gorm"
)

type repoBarberFavorit struct {
	Conn *gorm.DB
}

func NewRepoBarberFavorit(Conn *gorm.DB) ibarberfavorit.Repository {
	return &repoBarberFavorit{Conn}
}

func (db *repoBarberFavorit) GetDataBy(BarberId int, UserID int) (result *models.BarberFavorit, err error) {
	var (
		logger = logging.Logger{}
		data   = &models.BarberFavorit{}
	)
	query := db.Conn.Raw(`select * from barber_favorit where barber_id = ? and user_id = ? `, BarberId, UserID).Scan(&data) //Find(&result)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return result, models.ErrNotFound
		}
		return result, err
	}
	return data, nil
}
func (db *repoBarberFavorit) GetList(queryparam models.ParamListGeo) (result []*models.BarberFavoritList, err error) {
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
	sField := fmt.Sprintf(`
		b.barber_id,b.barber_cd,b.barber_name,
		b.address,b.latitude,b.longitude,
		b.operation_start,b.operation_end,
		b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type,true as is_favorit,
		fn_distance(%f,%f,b.latitude,b.longitude) as distance,
		(
			select (sum(fr.barber_rating)/count(fr.order_id))::float
			from feedback_rating fr 
			where fr.barber_id = b.barber_id 
		)::float as barber_rating
	`, queryparam.Latitude, queryparam.Longitude)

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(b.barber_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(b.barber_name) LIKE ?" //queryparam.Search
		}
		query = db.Conn.Table("barber b ").Select(sField).Joins(`
		join barber_favorit a
			on a.barber_id = b.barber_id 
		`).Joins(`
		left join sa_file_upload c
			on b.file_id = c.file_id 
	`).Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		query = db.Conn.Table("barber b ").Select(sField).Joins(`
		join barber_favorit a
			on a.barber_id = b.barber_id 
		`).Joins(`
		left join sa_file_upload c
			on b.file_id = c.file_id 
	`).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	}

	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (db *repoBarberFavorit) Create(data *models.BarberFavorit) error {
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
func (db *repoBarberFavorit) Update(ID int, data interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.OrderH{}).Where("order_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberFavorit) Delete(BarberId int, UserID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	// query := db.Conn.Where("order_id = ?", ID).Delete(&models.OrderH{})
	query := db.Conn.Exec("Delete From barber_favorit WHERE barber_id = ? and user_id = ?", BarberId, UserID)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarberFavorit) Count(queryparam models.ParamListGeo) (result int, err error) {
	var (
		sWhere  = ""
		logger  = logging.Logger{}
		query   *gorm.DB
		_result int64 = 0
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
	`).Where(sWhere, queryparam.Search).Count(&_result)
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
	`).Where(sWhere).Count(&_result)
	}
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return int(_result), nil
}
