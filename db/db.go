package db

import (

	// _ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var db *sql.DB
var db *gorm.DB
var err error

func Init() {
	// conf := config.GetConfig()
	another := os.Getenv("DB_URL")
	fmt.Println(another)
	// dsn := "host=localhost user=postgres password=root dbname=" + conf.DB_NAME + " port=" + conf.DB_PORT + ""
	db, err = gorm.Open(postgres.Open(another), &gorm.Config{})
	// connStr := conf.DB_USERNAME + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME
	// db, err = sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}

}

func CreateConn() *gorm.DB {
	return db.Debug() // change in prod
}
