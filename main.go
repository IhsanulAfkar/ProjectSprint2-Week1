package main

import (
	"week1/db"
	"week1/routes"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db.Init()
	r:= routes.Init()
	r.Run(":8080")
}