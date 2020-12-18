package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nocubicles/develytica/src/services"
	"github.com/nocubicles/develytica/src/utils"
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

	go services.ScanAndDoSyncs()

	log.Println("Listening..")
	router := router()

	PORT := os.Getenv("PORT")

	srv := fmt.Sprintf(":%s", PORT)

	err = http.ListenAndServe(srv, router)

	if err != nil {
		log.Fatal(err)
	}

}
