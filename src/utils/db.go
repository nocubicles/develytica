package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/nocubicles/skillbase.io/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error

	err = godotenv.Load(".env")

	if err != nil {
		panic("cannot load .env file")
	}
	dsn := os.Getenv("DBConnectionString")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		panic("failed to open db connection")
	}

	err = db.AutoMigrate(
		&models.Tenant{},
		&models.Sync{},
		&models.SyncHistory{},
		&models.User{},
		&models.Ad{},
		&models.Session{},
		&models.UserClaim{},
		&models.Sync{},
		&models.SyncHistory{},
		&models.GithubOrganization{},
		&models.GithubRepo{},
	)

	if err != nil {
		panic(err)
	}
}

//DbConnection returns the connection to use the db
func DbConnection() *gorm.DB {

	return db
}
