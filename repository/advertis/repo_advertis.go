package repoadvertis

import (
	"fmt"
	iadvertis "nuryanto2121/cukur_in_user/interface/advertis"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"

	"github.com/jinzhu/gorm"
)

type repoAdvertis struct {
	Conn *gorm.DB
}

func NewRepoAdvertis(Conn *gorm.DB) iadvertis.Repository {
	return &repoAdvertis{Conn}
}

func (db *repoAdvertis) GetDataBy(ID int) (result *models.Advertis, err error) {
	var (
		logger    = logging.Logger{}
		mAdvertis = &models.Advertis{}
	)
	query := db.Conn.Where("advertis_id = ? ", ID).Find(mAdvertis)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mAdvertis, nil
}

func (db *repoAdvertis) GetList(queryparam models.ParamList) (result []*models.ListAdvertis, err error) {

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
			select 	a.advertis_id ,		a.title ,
			a.descs ,			a.advertis_status ,
			a.slide_duration ,	a.start_date ,
			a.end_date ,		a.file_id ,
			fu.file_name ,		fu.file_path ,
			fu.file_type 
		from advertis a 
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

func (db *repoAdvertis) Create(data *models.Advertis) error {
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
func (db *repoAdvertis) Update(ID int, data map[string]interface{}) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Model(models.Advertis{}).Where("advertis_id = ?", ID).Updates(data)
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoAdvertis) Delete(ID int) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Where("advertis_id = ?", ID).Delete(&models.Advertis{})
	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}

func (db *repoAdvertis) Count(queryparam models.ParamList) (result int, err error) {
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
			sWhere += " and (lower(title) LIKE ? )" //+ queryparam.Search
		} else {
			sWhere += "(lower(title) LIKE ? )" //queryparam.Search
		}
		query = db.Conn.Model(&models.Advertis{}).Where(sWhere, queryparam.Search).Count(&result)
	} else {
		query = db.Conn.Model(&models.Advertis{}).Where(sWhere).Count(&result)
	}
	// end where

	logger.Query(fmt.Sprintf("%v", query.QueryExpr())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return result, nil
}
