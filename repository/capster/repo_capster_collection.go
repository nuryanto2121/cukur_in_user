package repocapstercollection

import (
	"fmt"
	icapster "nuryanto2121/cukur_in_user/interface/capster"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"

	"gorm.io/gorm"
)

type repoCapsterCollection struct {
	Conn *gorm.DB
}

func NewRepoCapsterCollection(Conn *gorm.DB) icapster.Repository {
	return &repoCapsterCollection{Conn}
}

func (db *repoCapsterCollection) GetDataBy(ID int, GeoBarber models.GeoBarber) (result *models.CapsterList, err error) {
	var (
		logger       = logging.Logger{}
		mCapsterList = &models.CapsterList{}
	)
	sSql := fmt.Sprintf(`
		SELECT * FROM fbarber_capster_s(%f,%f)
		WHERE capster_id = ?
	`, GeoBarber.Latitude, GeoBarber.Longitude)
	query := db.Conn.Raw(sSql, ID).First(&mCapsterList)
	// query := db.Conn.Where("capster_id = ? ", ID).Find(mCapsterList)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mCapsterList, nil
}
func (db *repoCapsterCollection) GetListFileCapter(ID int) (result []*models.SaFileOutput, err error) {
	var (
		logger = logging.Logger{}
	)
	query := db.Conn.Table("capster_collection").Select(`
			capster_collection.file_id,sa_file_upload.file_name,
			sa_file_upload.file_path, sa_file_upload.file_type
	`).Joins(`
		Inner Join sa_file_upload ON sa_file_upload.file_id = capster_collection.file_id
	`).Where("capster_id = ?", ID).Find(&result)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return result, models.ErrNotFound
		}
		return nil, err
	}
	return result, nil
}
func (db *repoCapsterCollection) GetList(queryparam models.ParamListGeo) (result []*models.CapsterList, err error) {

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
			sWhere += " and lower(name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(name) LIKE ?" //queryparam.Search
		}

		query = db.Conn.Table("v_capster").Select(`
				capster_id,user_name,name,
				is_active,file_id,file_name,
				file_path,file_type,rating,
				user_type,user_input,time_edit,
				in_use
		`).Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)
	} else {
		query = db.Conn.Table("v_capster").Select(`
			capster_id,user_name,name,
			is_active,file_id,file_name,
			file_path,file_type,rating,
			user_type,user_input,time_edit,
			in_use
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
func (db *repoCapsterCollection) Count(queryparam models.ParamListGeo) (result int, err error) {
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
			sWhere += " and lower(name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(name) LIKE ?" //queryparam.Search
		}
		query = db.Conn.Table("v_capster").Select(`
			v_capster.capster_id,v_capster.name,v_capster.is_active, 0 as rating
		`).Where(sWhere, queryparam.Search).Count(&_result)
	} else {
		query = db.Conn.Table("v_capster").Select(`
			v_capster.capster_id,v_capster.name,v_capster.is_active, 0 as rating
		`).Where(sWhere).Count(&_result)
	}

	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return int(_result), nil
}
