package postgresdb

import (
	"fmt"
	"log"
	version "nuryanto2121/cukur_in_user/middleware/versioning"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/setting"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // add database driver bridge
)

var Conn *gorm.DB

func Setup() {
	now := time.Now()
	var err error
	fmt.Print(setting.FileConfigSetting.Database)
	connectionstring := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		setting.FileConfigSetting.Database.User,
		setting.FileConfigSetting.Database.Password,
		setting.FileConfigSetting.Database.Name,
		setting.FileConfigSetting.Database.Host,
		setting.FileConfigSetting.Database.Port)
	fmt.Printf("%s", connectionstring)
	Conn, err = gorm.Open(setting.FileConfigSetting.Database.Type, connectionstring)
	if err != nil {
		log.Printf("connection.setup err : %v", err)
		panic(err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.FileConfigSetting.Database.TablePrefix + defaultTableName
	}
	Conn.SingularTable(true)
	Conn.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	Conn.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	Conn.Callback().Delete().Replace("gorm:delete", deleteCallback)

	Conn.DB().SetMaxIdleConns(10)
	Conn.DB().SetMaxOpenConns(100)

	go autoMigrate()

	timeSpent := time.Since(now)
	log.Printf("Config database is ready in %v", timeSpent)
}

// autoMigrate : create or alter table from struct
func autoMigrate() {
	// Add auto migrate bellow this line
	Conn.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE OR REPLACE FUNCTION fn_distance(lat1 FLOAT, lon1 FLOAT, lat2 FLOAT, lon2 FLOAT) RETURNS FLOAT AS $$
		DECLARE                                                   
			x float = 111.12  * (lat2 - lat1);                           
			y float = 111.12  * (lon2 - lon1) * cos(lat1 / 92.215);        
		BEGIN                                                     
			RETURN sqrt(x * x + y * y);                               
		END  
		$$ LANGUAGE plpgsql;

		CREATE OR replace FUNCTION public.fbarber_beranda_user_s(p_latitude FLOAT, p_longitude FLOAT)
		RETURNS TABLE(
			barber_id integer, 	barber_cd varchar, 	barber_name varchar, 
			address varchar, 	latitude float, 	longitude float,
			operation_start timestamp, 	operation_end timestamp,
			is_active bool, 	file_id integer, 	file_name varchar, 	file_path varchar, 	file_type varchar, 
			is_favorit bool, 	distance float,		barber_rating float, is_barber_open bool
		)
		LANGUAGE plpgsql
		AS $function$
			DECLARE v_id INTEGER; 
			BEGIN 	
				RETURN QUERY                
					select
								b.barber_id,b.barber_cd,b.barber_name,
								b.address,b.latitude,b.longitude,
								b.operation_start,b.operation_end,
								b.is_active,c.file_id ,c.file_name ,c.file_path ,c.file_type,
								case when a.barber_id  is not null then true else false end as is_favorit,
								fn_distance(p_latitude,p_longitude,b.latitude,b.longitude) as distance,
								(
									select (sum(fr.barber_rating)/count(fr.order_id))::float
									from feedback_rating fr 
									where fr.barber_id = b.barber_id 
								)::float as barber_rating,
								(
									case when now() between (now()::date + b.operation_start::time) and (now()::date + b.operation_end ::time) then true else false end
								)::bool as is_barber_open
						from barber b 
						left join barber_favorit a
							on a.barber_id = b.barber_id 
						left join sa_file_upload c
								on b.file_id = c.file_id 
			;
						
			END;
			$function$
		;

		CREATE OR replace FUNCTION public.fbarber_capster_s(p_latitude FLOAT, p_longitude FLOAT)
		RETURNS TABLE(
			capster_id integer, 	capster_name varchar, 			is_active bool, 
			file_id integer, 		file_name varchar, 		file_path varchar, 			file_type varchar,
			barber_id integer, 		barber_name varchar, 	distance float,
			capster_rating float,  	is_barber_open bool,	operation_start timestamp, 	operation_end timestamp,
			is_barber_active bool,	join_date timestamp,	barber_rating float
		)
		LANGUAGE plpgsql
	   AS $function$
			DECLARE v_id INTEGER; 
			BEGIN 	
				  RETURN QUERY                
						  select
						   a.user_id as capster_id, 	a.name as capster_name, 		a.is_active,
						   d.file_id,	d.file_name,	d.file_path,	d.file_type, 
						   ab.barber_id , b.barber_name ,
						   fn_distance(p_latitude,p_longitude,b.latitude,b.longitude) as distance,
						   (
							   select (sum(fr.capster_rating)/count(fr.order_id))::float
							   from feedback_rating fr 
							   where fr.capster_id = a.user_id 
						   )::float as capster_rating,
						   (
							   case when now()::timestamp without time zone between (now()::date + b.operation_start::time)::timestamp without time zone and (now()::date + b.operation_end ::time)::timestamp without time zone then true else false end
						   )::bool as is_barber_open,
						   b.operation_start ,b.operation_end ,b.is_active as is_barber_active,
						   a.join_date ,
						   (
								   select (sum(fr.barber_rating)/count(fr.order_id))::float
								   from feedback_rating fr 
								   where fr.barber_id = b.barber_id 
							   )::float as barber_rating
					   from ss_user a
					   inner join barber_capster ab
							   ON ab.capster_id = a.user_id	
					   inner join barber b
							   on b.barber_id = ab.barber_id 
					   left join sa_file_upload d
							   ON d.file_id = a.file_id
			;
					   
			END;
			$function$
	   ;

	`)
	log.Println("STARTING AUTO MIGRATE ")
	Conn.AutoMigrate(
		version.SsVersion{},
		models.BarberFavorit{},
		models.FeedbackRating{},
	)

	log.Println("FINISHING AUTO MIGRATE ")
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		// nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("TimeInput"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(util.GetTimeNow())
			}
		}

		if modifyTimeField, ok := scope.FieldByName("TimeEdit"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(util.GetTimeNow())
			}
		}

	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("TimeEdit", util.GetTimeNow())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
