package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nocubicles/skillbase.io/src/utils"
)

func init() {
	if os.Getenv("GO_ENV") != "PRODUCTION" {
		err := godotenv.Load(".env")

		if err != nil {
			panic("cannot load .env file")
		}

	}

}

func main() {
	db := utils.DbConnection()
	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	defer sqlDB.Close()

	log.Println("Listening..")
	router := router()
	err = http.ListenAndServe(":3000", router)

	if err != nil {
		log.Fatal(err)
	}
}
