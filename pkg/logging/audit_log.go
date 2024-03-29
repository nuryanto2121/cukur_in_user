package logging

import (
	"context"
	"fmt"
	"log"
	"nuryanto2121/cukur_in_user/pkg/monggodb"
	util "nuryanto2121/cukur_in_user/pkg/utils"
	"time"
)

type auditLog struct {
	ID        int64     `bson:"id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Level     string    `json:"level"`
	UUID      string    `json:"uuid"`
	FuncName  string    `json:"func_name"`
	FileName  string    `json:"file_name"`
	Line      int       `json:"line"`
	Time      string    `json:"time"`
	Message   string    `json:"message"`
}

func (a *auditLog) saveAudit() {

	a.ID = util.GetTimeNow().Unix()
	a.Message = "API User : " + a.Message
	result, err := monggodb.MCon.Collection("auditlogs").InsertOne(context.TODO(), a)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Inserted a single document: ", result.InsertedID)

}
