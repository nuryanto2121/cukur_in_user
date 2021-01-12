package postgresgorm

import (
	"fmt"
	"log"
	version "nuryanto2121/cukur_in_user/middleware/versioning"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/setting"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Conn *gorm.DB

func Setup() {
	now := time.Now()
	var err error

	connectionstring := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		setting.FileConfigSetting.Database.Host,
		setting.FileConfigSetting.Database.User,
		setting.FileConfigSetting.Database.Password,
		setting.FileConfigSetting.Database.Name,
		setting.FileConfigSetting.Database.Port)
	fmt.Printf("%s", connectionstring)

	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Jakarta"
	Conn, err = gorm.Open(postgres.Open(connectionstring), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   setting.FileConfigSetting.Database.TablePrefix,
			SingularTable: true,
		},
	})

	if err != nil {
		log.Printf("connection.setup err : %v", err)
		panic(err)
	}

	// Conn.SingularTable(true)
	// Conn.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	// Conn.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	// Conn.Callback().Delete().Replace("gorm:delete", deleteCallback)

	sqlDB, err := Conn.DB()
	if err != nil {
		log.Printf("connection.setup DB err : %v", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	// sqlDB.Sing
	// Conn.DB().SetMaxIdleConns(10)
	// Conn.DB().SetMaxOpenConns(100)

	go autoMigrate()

	timeSpent := time.Since(now)
	log.Printf("Config database is ready in %v", timeSpent)
}

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

		CREATE OR replace FUNCTION public.fbarber_beranda_user_s(p_latitude FLOAT, p_longitude FLOAT,p_user_id integer)
		RETURNS TABLE(
			barber_id integer, 			barber_cd varchar, 				barber_name varchar, 
			address varchar, 			latitude float, 				longitude float,
			operation_start timestamp, 	operation_end timestamp, 		is_active bool, 	
			file_id integer, 			file_name varchar, 				file_path varchar, 	
			file_type varchar,			is_favorit bool, 				distance float,					
			barber_rating float, 		is_barber_open bool,			total_user_order integer
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
								)::bool as is_barber_open,
								(
									select count(fr.user_id)
									from feedback_rating fr 
									where fr.barber_id = b.barber_id 
								)::integer as total_user_order
						from barber b 
						left join barber_favorit a
							on a.barber_id = b.barber_id 
							and a.user_id = p_user_id
						left join sa_file_upload c
								on b.file_id = c.file_id 
			;
						
			END;
			$function$
		;

		CREATE OR replace  FUNCTION public.fbarber_capster_s(p_latitude FLOAT, p_longitude FLOAT)
		RETURNS TABLE(
			capster_id integer, 	capster_name varchar, 			is_active bool, 
			file_id integer, 		file_name varchar, 		file_path varchar, 			file_type varchar,
			barber_id integer, 		barber_name varchar, 	distance float,
			capster_rating float,  	is_barber_open bool,	operation_start timestamp, 	operation_end timestamp,
			is_barber_active bool,	join_date timestamp,	barber_rating float, 		is_busy bool,
			length_of_work varchar,	latitude float8,		longitude float8
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
								)::float as barber_rating,
							(
							case when
								(select count(oh.status)
								from order_h oh 
								where oh.capster_id = a.user_id 
								and oh.order_date::date = now()::date
								and oh.status = 'P') = 0 then false else true end
							) as is_busy,
							( 	
								case when extract(year from age(current_date,a.join_date)) > 0 
							 		then (TO_CHAR(age(current_date, a.join_date), 'YY')::integer)::varchar ||' Tahun'
							 	when extract(month from age(current_date,a.join_date)) > 0
							 		then (TO_CHAR(age(current_date, a.join_date), 'mm')::integer)::varchar ||' Bulan'
							 	else (TO_CHAR(age(current_date, a.join_date), 'DD')::integer)::varchar ||' Hari'
							  	end 
							)::varchar as length_of_work,
							b.latitude ,b.longitude 
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
		// models.SsUser{},
		version.SsVersion{},
		models.BarberFavorit{},
		models.FeedbackRating{},
		models.BookingCapster{},
		models.Notification{},
		models.Advertise{},
	)

	log.Println("FINISHING AUTO MIGRATE ")
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(db *gorm.DB) {
	if db.Statement.Error == nil {
		TimeInput := db.Statement.Schema.LookUpField("TimeInput")
		TimeInput.Set(db.Statement.ReflectValue, util.GetTimeNow())

		TimeEdit := db.Statement.Schema.LookUpField("TimeEdit")
		TimeEdit.Set(db.Statement.ReflectValue, util.GetTimeNow())
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if db.Statement.Changed() {
		db.Statement.SetColumn("TimeEdit", time.Now())
	}

}
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
