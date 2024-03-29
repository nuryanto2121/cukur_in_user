package repobarber

import (
	"context"
	"fmt"
	ibarber "nuryanto2121/cukur_in_user/interface/barber"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/setting"
	"time"

	"gorm.io/gorm"
)

type repoBarber struct {
	Conn *gorm.DB
}

func NewRepoBarber(Conn *gorm.DB) ibarber.Repository {
	return &repoBarber{Conn}
}

func (db *repoBarber) GetDataBy(ID int) (result *models.Barber, err error) {
	var (
		logger  = logging.Logger{}
		mBarber = &models.Barber{}
	)
	query := db.Conn.Where("barber_id = ? ", ID).Find(mBarber)
	logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mBarber, nil
}

func (db *repoBarber) GetDataByList(ctx context.Context, ID int, UserID int, GeoBarber models.GeoBarber) (result *models.BarbersList, err error) {
	var (
		logger  = logging.Logger{}
		mBarber = &models.BarbersList{}
	)
	sSql := fmt.Sprintf(`
		SELECT * FROM fbarber_beranda_user_s(%f,%f,%d)
		WHERE barber_id = ?
	`, GeoBarber.Latitude, GeoBarber.Longitude, UserID)

	// query := db.Conn.Where("barber_id = ? ", ID).Find(mBarber)
	// query := db.Conn.Raw(sSql, ID).First(&mBarber)
	// query := db.Conn.WithContext(ctx).Raw(sSql, ID).First(&mBarber)
	query := db.Conn.Debug().Raw(sSql, ID).First(&mBarber)

	// query := db.Conn.Session(&Session{DryRun: true}).Raw(sSql, ID).First(&mBarber)
	// query.SQL.String()
	// query.Vars
	logger.Query(fmt.Sprintf("%v", query))
	// logger.Query(fmt.Sprintf("%v", query.Vars))
	// logger.Query(fmt.Sprintf("%v", query.Vars))

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
		logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String()))
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
		logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String()))
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
func (db *repoBarber) GetDataBarber(BarberID int) (result *models.Barber, err error) {
	var (
		logger  = logging.Logger{}
		mBarber = &models.Barber{}
	)

	query := db.Conn.Where("barber_id = ? ", BarberID).Find(mBarber)
	logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String()))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mBarber, nil
}
func (db *repoBarber) GetListOld(UserID int, queryparam models.ParamListGeo) (result []*models.BarbersList, err error) {

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
					(now()::date + b.operation_start::time) as operation_start, (now()::date + b.operation_end ::time) as operation_end,
					b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type,
					case when a.barber_id  is not null then true else false end as is_favorit,
					fn_distance(%f,%f,b.latitude,b.longitude) as distance,
					coalesce ((
						select (sum(fr.barber_rating)/count(fr.order_id))::float
						from feedback_rating fr 
						where fr.barber_id = b.barber_id 
					),0)::float as barber_rating,
					(
						case when now() between (now()::date + b.operation_start::time) and (now()::date + b.operation_end ::time) then 1 else 0 end
					)::boolean as is_barber_open,								
					(
						select count(fr.user_id)
						from feedback_rating fr 
						where fr.barber_id = b.barber_id 
					)::integer as total_user_order
	
	`, queryparam.Latitude, queryparam.Longitude)

	sJoin := fmt.Sprintf(`
		left join barber_favorit a
		on a.barber_id = b.barber_id 
		and a.user_id = %d
	`, UserID)

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(b.barber_name) iLIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(b.barber_name) iLIKE ?" //queryparam.Search
		}
		query = db.Conn.Table("barber b ").Select(sField).Joins(sJoin).Joins(`
		left join sa_file_upload c
		on b.file_id = c.file_id 
	`).Where(sWhere, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

	} else {
		query = db.Conn.Table("barber b").Select(sField).Joins(sJoin).Joins(`
		left join sa_file_upload c
		on b.file_id = c.file_id 
		`).Where(sWhere).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

	}

	logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String())) //cath to log query string
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	fmt.Printf("%v", len(result))
	return result, nil
}
func (db *repoBarber) GetList(UserID int, queryparam models.ParamListGeo) (result []*models.BarbersList, err error) {

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

	sSql := fmt.Sprintf(` select * from (
		select
								b.barber_id,b.barber_cd,b.barber_name,
								b.address,b.latitude,b.longitude,
								(now()::date + b.operation_start::time) as operation_start, (now()::date + b.operation_end ::time) as operation_end,
								b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type,
								case when a.barber_id  is not null then true else false end as is_favorit,
								fn_distance(%f,%f,b.latitude,b.longitude) as distance,
								coalesce ((
									select (sum(fr.barber_rating)/count(fr.order_id))::float
									from feedback_rating fr 
									where fr.barber_id = b.barber_id 
								),0)::float as barber_rating,
								(
									case when now() between (now()::date + b.operation_start::time) and (now()::date + b.operation_end ::time) then true else false end
								)::boolean as is_barber_open,								
								(
									select count(fr.user_id)
									from feedback_rating fr 
									where fr.barber_id = b.barber_id 
								)::integer as total_user_order
						from barber b 
						left join barber_favorit a
							on a.barber_id = b.barber_id 
							and a.user_id = %d
						left join sa_file_upload c
								on b.file_id = c.file_id 
						) a 
	`, queryparam.Latitude, queryparam.Longitude, UserID)

	if queryparam.Search != "" {
		if sWhere != "" {
			sWhere += " and lower(barber_name) LIKE ?" //+ queryparam.Search
		} else {
			sWhere += "lower(barber_name) LIKE ?" //queryparam.Search
		}
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		if orderBy != "" {
			sSql += fmt.Sprintf("\n order by %s", orderBy)
		}
		sSql += fmt.Sprintf("\n offset %d limit %d", pageNum, pageSize)
		query = db.Conn.Raw(sSql, queryparam.Search).Find(&result)
		// sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		// fmt.Println(sSql)
		// query = db.Conn.Raw(sSql, queryparam.Search).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

	} else {
		sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		if orderBy != "" {
			sSql += fmt.Sprintf("\n order by %s", orderBy)
		}
		sSql += fmt.Sprintf("\n offset %d limit %d", pageNum, pageSize)
		query = db.Conn.Raw(sSql).Find(&result)
		// sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
		// fmt.Println(sSql)
		// query = db.Conn.Raw(sSql).Offset(pageNum).Limit(pageSize).Order(orderBy).Find(&result)

	}

	logger.Query(fmt.Sprintf("%v", sSql)) //cath to log query string
	err = query.Error
	fmt.Printf("%v", len(result))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}
func (db *repoBarber) GetScheduleTime(BarberID int) (result interface{}, err error) {
	type Results struct {
		ScheduleTime time.Time `json:"schedule_time"`
	}

	var (
		logger = logging.Logger{}
		op     = []*Results{}
		query  *gorm.DB
	)
	result = 0

	sSql := fmt.Sprintf(` SELECT h.schedule_time
	FROM   (select
				generate_series (b2.operation_start::timestamp 
						  , b2.operation_end::timestamp 
						  , interval '30m') as schedule_time,
				   b2.operation_start ,
				   b2.operation_end 
			 from barber b2 
			 where b2.barber_id = %d
			)h
	WHERE  EXTRACT(ISODOW FROM h.schedule_time) < 6
	and h.schedule_time::time >= now()::time
	AND   now()::time >= h.operation_start::time 
	AND   now()::time <= h.operation_end::time ;
	`, BarberID)

	// sSql = fmt.Sprintf(sSql+` WHERE %s`, sWhere)
	query = db.Conn.Raw(sSql).Find(&op)

	// end where

	logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return op, nil
}
func (db *repoBarber) Create(data *models.Barber) error {
	var (
		logger = logging.Logger{}
		err    error
	)
	query := db.Conn.Create(data)
	logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String())) //cath to log query string
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
	logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String())) //cath to log query string
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
	logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String())) //cath to log query string
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
func (db *repoBarber) CountOld(UserID int, queryparam models.ParamListGeo) (result int, err error) {

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
								(now()::date + b.operation_start::time) as operation_start, (now()::date + b.operation_end ::time) as operation_end,
								b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type,
								case when a.barber_id  is not null then true else false end as is_favorit,
								fn_distance(?,?,b.latitude,b.longitude) as distance,
								coalesce ((
									select (sum(fr.barber_rating)/count(fr.order_id))::float
									from feedback_rating fr 
									where fr.barber_id = b.barber_id 
								),0)::float as barber_rating,
								(
									case when now() between (now()::date + b.operation_start::time) and (now()::date + b.operation_end ::time) then 1 else 0 end
								)::boolean as is_barber_open,								
								(
									select count(fr.user_id)
									from feedback_rating fr 
									where fr.barber_id = b.barber_id 
								)::integer as total_user_order
		`, queryparam.Latitude, queryparam.Longitude).Joins(`
			left join barber_favorit a
				on a.barber_id = b.barber_id 
				and a.user_id = ?
		`, UserID).Joins(`
		left join sa_file_upload c
		on b.file_id = c.file_id 
	`).Where(sWhere, queryparam.Search).Count(&_result)
		//query = db.Conn.Model(&models.Barber{}).Where(sWhere, queryparam.Search).Count(&_result)
	} else {
		query = db.Conn.Table("barber b ").Select(`
								b.barber_id,b.barber_cd,b.barber_name,
								b.address,b.latitude,b.longitude,
								(now()::date + b.operation_start::time) as operation_start, (now()::date + b.operation_end ::time) as operation_end,
								b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type,
								case when a.barber_id  is not null then true else false end as is_favorit,
								fn_distance(?,?,b.latitude,b.longitude) as distance,
								coalesce ((
									select (sum(fr.barber_rating)/count(fr.order_id))::float
									from feedback_rating fr 
									where fr.barber_id = b.barber_id 
								),0)::float as barber_rating,
								(
									case when now() between (now()::date + b.operation_start::time) and (now()::date + b.operation_end ::time) then 1 else 0 end
								)::boolean as is_barber_open,								
								(
									select count(fr.user_id)
									from feedback_rating fr 
									where fr.barber_id = b.barber_id 
								)::integer as total_user_order
		`, queryparam.Latitude, queryparam.Longitude).Joins(`
			left join barber_favorit a
				on a.barber_id = b.barber_id 
				and a.user_id = ?
		`, UserID).Joins(`
		left join sa_file_upload c
		on b.file_id = c.file_id 
	`).Where(sWhere).Count(&_result)
		// query = db.Conn.Model(&models.Barber{}).Where(sWhere).Count(&result)
	}
	// end where

	logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return int(_result), nil
}
func (db *repoBarber) Count(UserID int, queryparam models.ParamListGeo) (result int, err error) {

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
	sSql := fmt.Sprintf(` select count(*) as cnt from (
		select
								b.barber_id,b.barber_cd,b.barber_name,
								b.address,b.latitude,b.longitude,
								(now()::date + b.operation_start::time) as operation_start, (now()::date + b.operation_end ::time) as operation_end,
								b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type,
								case when a.barber_id  is not null then true else false end as is_favorit,
								fn_distance(%f,%f,b.latitude,b.longitude) as distance,
								(
									select (sum(fr.barber_rating)/count(fr.order_id))::float
									from feedback_rating fr 
									where fr.barber_id = b.barber_id 
								)::float as barber_rating,
								(
									case when now() between (now()::date + b.operation_start::time) and (now()::date + b.operation_end ::time) then 1 else 0 end
								)::boolean as is_barber_open,								
								(
									select count(fr.user_id)
									from feedback_rating fr 
									where fr.barber_id = b.barber_id 
								)::integer as total_user_order
						from barber b 
						left join barber_favorit a
							on a.barber_id = b.barber_id 
							and a.user_id = %d
						left join sa_file_upload c
								on b.file_id = c.file_id 
						) a
	`, queryparam.Latitude, queryparam.Longitude, UserID)

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

	logger.Query(fmt.Sprintf("%v", query.Statement.SQL.String())) //cath to log query string
	err = query.Error
	if err != nil {
		return 0, err
	}

	return op.Cnt, nil
}
