// package db

// import (

// 	"context"
// 	"os"

// 	"github.com/jackc/pgx/v5"
	
// )

// // var db *sql.DB
// var db *pgx.Conn
// var err error

// func Init() {
// 	another := os.Getenv("DB_URL")
// 	db, err = gorm.Open(postgres.Open(another), &gorm.Config{})

// 	if err != nil {
// 		panic(err.Error())
// 	}

// }

// func CreateConn() *pgx.Conn {
// 	return db
// }
