package repoadvertis

import (
	"fmt"
	iadvertis "nuryanto2121/cukur_in_user/interface/advertise"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"

	"gorm.io/gorm"
)

type repoAdvertise struct {
	Conn *gorm.DB
}

func NewRepoAdvertise(Conn *gorm.DB) iadvertis.Repository {
	return &repoAdvertise{Conn}
}

func (db *repoAdvertise) GetDataBy(ID int) (result *models.Advertise, err error) {
	var (
		logger     = logging.Logger{}
		mAdvertise = &models.Advertise{}
	)
	query := db.Conn.Where("advertise_id = ? ", ID).Find(mAdvertise)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mAdvertise, nil
}

func (db *repoAdvertise) GetList(queryparam models.ParamList) (result []*models.ListAdvertise, err error) {

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

	sSql := `
		SELECT * FROM (
			select 	a.advertise_id ,		a.title ,
			a.descs ,			a.advertise_status ,
			a.slide_duration ,	a.start_date ,
			a.end_date ,		a.file_id ,
			fu.file_name ,		fu.file_path ,
			fu.file_type 
		from advertise a 
		join sa_file_upload fu 
			on a.file_id = fu.file_id 
		) a
	
	`
	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and (lower(title) LIKE ?)"
		} else {
			sWhere += "(lower(title) LIKE ?)"
		}
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		query = db.Conn.Raw(sSql, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		query = db.Conn.Raw(sSql).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
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

func (db *repoAdvertise) Create(data *models.Advertise) error {
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
func (db *repoAdvertise) Update(ID int, data map[string]interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.Advertise{}).Where("advertise_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoAdvertise) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where("advertise_id = ?", ID).Delete(&models.Advertise{})
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoAdvertise) Count(queryparam models.ParamList) (result int, err error) {
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
			sWhere += " and (lower(title) LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower(title) LIKE ? )" //queryparam.Search
		}
		query = db.Conn.Model(&models.Advertise{}).Where(sWhere, queryparam.Search).Count(&_result)
	} else {
		query = db.Conn.Model(&models.Advertise{}).Where(sWhere).Count(&_result)
	}
	// end where

	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}
	result = int(_result)
	return result, nil
}
