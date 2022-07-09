package main

import (
	"hermes/db"
	"os"
	"strconv"
)

func main() {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		os.Exit(1)
	}

	app := App{}
	app.Initialize(db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     int32(dbPort),
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Catalog:  os.Getenv("DB_NAME"),
	})
	app.Run(":8024")
}
