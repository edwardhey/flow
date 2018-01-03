package models

import (
	"fmt"
	"time"

	"github.com/edwardhey/lib/datasource"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type UUID uint64

func (u UUID) String() string {
	return fmt.Sprintf("%d", u)
}

func init() {
	dns := "root:123456@tcp(mysql.localhost:3306)/flow"
	datasource.InitMysql(dns)
	db = datasource.GetMysql()

	objects := []interface{}{
		&Flow{},
		&Node{},
		&Line{},
	}

	db.AutoMigrate(objects...)

}

type Base struct {
	CTime time.Time `gorm:"column:ctime"`
	Mtime time.Time `gorm:"column:mtime"`
}
