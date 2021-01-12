package repofunction

import (
	"fmt"
	"nuryanto2121/cukur_in_user/models"
	"nuryanto2121/cukur_in_user/pkg/logging"
	"nuryanto2121/cukur_in_user/pkg/postgresdb"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type FN struct {
	Claims util.Claims
}

func (fn *FN) GenTransactionNo(BarberCd string) (string, error) {
	var (
		result string
		conn   *gorm.DB
		logger = logging.Logger{}
		mSeqNo = &models.SsSequenceNo{}
		t      = time.Now()
		year   = t.Year()
		month  = int(t.Month())
		sYear  = strconv.Itoa(year)
		sMonth = strconv.Itoa(month)

		// abjad  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	)
	if len(sMonth) == 1 {
		sMonth = fmt.Sprintf("0%s", sMonth)
	}
	if len(sYear) == 4 {
		sYear = sYear[len(sYear)-2:]
	}
	pref := fmt.Sprintf("%s%s", sMonth, sYear)
	conn = postgresdb.Conn

	query := conn.Where("sequence_cd = ?", BarberCd).Find(mSeqNo)
	logger.Query(fmt.Sprintf("%v", query))
	err := query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			mSeqNo.Prefix = fmt.Sprintf("%s", pref)
			mSeqNo.SeqNo = 1
			mSeqNo.SequenceCd = BarberCd
			mSeqNo.UserInput = fn.Claims.UserID
			mSeqNo.UserEdit = fn.Claims.UserID
			queryC := conn.Create(mSeqNo)
			logger.Query(fmt.Sprintf("%v", queryC.QueryExpr()))
			err = queryC.Error
			if err != nil {
				return "", err
			}
			result = fmt.Sprintf("%s/%s/0001", BarberCd, pref)
			return result, nil
		}
		return "", err
	}
	seq_no := ""

	if mSeqNo.Prefix == pref {
		mSeqNo.SeqNo += 1
		seq_no = strconv.Itoa(10000 + mSeqNo.SeqNo)
		seq_no = seq_no[len(seq_no)-4:]
	} else {
		mSeqNo.Prefix = fmt.Sprintf("%s", pref)
		mSeqNo.SeqNo = 1
		seq_no = "0001"
	}
	query = conn.Model(models.SsSequenceNo{}).Where("sequence_id = ?", mSeqNo.SequenceID).Updates(mSeqNo)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		return "", err
	}
	result = fmt.Sprintf("%s/%s/%s", BarberCd, pref, seq_no)

	return result, nil
}

func (fn *FN) GenBarberCode() (string, error) {
	var (
		result string
		conn   *gorm.DB
		logger = logging.Logger{}
		mSeqNo = &models.SsSequenceNo{}
		abjad  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	)
	conn = postgresdb.Conn

	prefixArr := strings.Split(abjad, "")
	fmt.Printf("%v", prefixArr)
	// ss := prefixArr[0]
	// query := conn.Table("barber").Select("max(barber_cd)") //
	query := conn.Where("sequence_cd = ?", "seq_barber").Find(mSeqNo)
	logger.Query(fmt.Sprintf("%v", query))
	err := query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			mSeqNo.Prefix = "AA"
			mSeqNo.SeqNo = 1
			mSeqNo.SequenceCd = "seq_barber"
			mSeqNo.UserInput = fn.Claims.UserID
			mSeqNo.UserEdit = fn.Claims.UserID
			queryC := conn.Create(mSeqNo)
			logger.Query(fmt.Sprintf("%v", queryC.QueryExpr()))
			err = queryC.Error
			if err != nil {
				return "", err
			}
			result = "AA01"
			return result, nil
		}
		return "", err
	}

	if mSeqNo.SeqNo == 99 {
		mSeqNo.SeqNo = 1

		pref := strings.Split(mSeqNo.Prefix, "")
		i1 := strings.Index(abjad, pref[0])
		i2 := strings.Index(abjad, pref[1])
		aa := len(prefixArr) - 1
		if i2 == aa {
			i2 = 0
			i1 += 1
		} else {
			i2 += 1
		}

		mSeqNo.Prefix = fmt.Sprintf("%s%s", prefixArr[i1], prefixArr[i2])

	} else {
		mSeqNo.SeqNo += 1
	}

	sNo := mSeqNo.SeqNo + 100

	runes := []rune(strconv.Itoa(sNo))
	no := string(runes[1:])
	query = conn.Model(models.SsSequenceNo{}).Where("sequence_id = ?", mSeqNo.SequenceID).Updates(mSeqNo)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		return "", err
	}
	result = fmt.Sprintf("%s%s", mSeqNo.Prefix, no)

	return result, nil
}

func (fn *FN) GetCountTrxCapster(ID int) int {
	var (
		logger = logging.Logger{}
		result = 0
		conn   *gorm.DB
	)
	OwnerID, _ := strconv.Atoi(fn.Claims.UserID)

	conn = postgresdb.Conn
	query := conn.Table("order_h oh").Select(`
		oh.order_id ,oh.order_no ,oh.barber_id ,b.barber_name ,b.owner_id ,oh.order_date ,oh.status ,oh.capster_id
	`).Joins(`
		join barber b on b.barber_id = oh.barber_id 
	`).Where(`
	oh.status in('P','N') AND b.owner_id = ? AND oh.capster_id = ?
		`, OwnerID, ID).Count(&result)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err := query.Error
	if err != nil {
		if err == models.ErrNotFound {
			return 0
		}
		logger.Error(err)
	}
	return result
}

func (fn *FN) GetCountTrxBarber(ID int) int {

	var (
		logger = logging.Logger{}
		result = 0
		conn   *gorm.DB
	)
	OwnerID, _ := strconv.Atoi(fn.Claims.UserID)

	conn = postgresdb.Conn
	query := conn.Table("order_h oh").Select(`
		oh.order_id ,oh.order_no ,oh.barber_id ,b.barber_name ,b.owner_id ,oh.order_date ,oh.status ,oh.capster_id
	`).Joins(`
		join barber b on b.barber_id = oh.barber_id 
	`).Where(`
	oh.status in('P','N') AND b.owner_id = ? AND oh.barber_id = ?
		`, OwnerID, ID).Count(&result)
	logger.Query(fmt.Sprintf("%v", query)) //cath to log query string
	err := query.Error
	if err != nil {
		if err == models.ErrNotFound {
			return 0
		}
		logger.Error(err)
	}
	return result
}

func (fn *FN) GetUserData() (result *models.SsUser, err error) {
	var (
		logger   = logging.Logger{}
		mCapster = &models.SsUser{}
		conn     *gorm.DB
	)
	conn = postgresdb.Conn
	query := conn.Where("user_id = ? ", fn.Claims.UserID).Find(mCapster)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mCapster, nil
}

func (fn *FN) GetBarberData(BarberID int) (result *models.Barber, err error) {
	var (
		logger  = logging.Logger{}
		mBarber = &models.Barber{}
		conn    *gorm.DB
	)
	conn = postgresdb.Conn
	query := conn.Where("barber_id = ? ", BarberID).Find(mBarber)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mBarber, nil
}
func (fn *FN) GetCapsterData(CapsterID int) (result *models.SsUser, err error) {
	var (
		logger   = logging.Logger{}
		mCapster = &models.SsUser{}
		conn     *gorm.DB
	)
	conn = postgresdb.Conn
	query := conn.Where("user_id = ? ", CapsterID).Find(mCapster)
	logger.Query(fmt.Sprintf("%v", query))
	err = query.Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return mCapster, nil
}
func (fn *FN) InTimeActiveBarber(data *models.Barber, orderDate time.Time) bool {
	var (
		P = fmt.Println
	)
	// timeStart := data.OperationStart.Date(timeOrder.Year(), timeOrder.Month(), timeOrder.Day())
	timeStart := time.Date(orderDate.Year(), orderDate.Month(), orderDate.Day(), data.OperationStart.Hour(),
		data.OperationStart.Minute(), data.OperationStart.Second(), data.OperationStart.Nanosecond(), data.OperationStart.Local().Location())

	timeEnd := time.Date(orderDate.Year(), orderDate.Month(), orderDate.Day(), data.OperationEnd.Hour(),
		data.OperationEnd.Minute(), data.OperationEnd.Second(), data.OperationEnd.Nanosecond(), data.OperationEnd.Local().Location())

	timeOrder := time.Date(orderDate.Year(), orderDate.Month(), orderDate.Day(), orderDate.Hour(),
		orderDate.Minute(), orderDate.Second(), orderDate.Nanosecond(), orderDate.Local().Location())

	P(timeStart)
	P(timeEnd)
	P(timeOrder)

	if timeOrder.Before(timeEnd) && timeOrder.After(timeStart) {
		return true
	} else {
		return false
	}

}
