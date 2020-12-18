package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/nocubicles/develytica/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	var err error

	if os.Getenv("GO_ENV") != "PRODUCTION" {
		err := godotenv.Load(".env")

		if err != nil {
			panic("cannot load .env file")
		}
	}
	dsn := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

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
		&models.Organization{},
		&models.Repo{},
		&models.Assignee{},
		&models.Issue{},
		&models.Label{},
		&models.IssueAssignee{},
		&models.IssueLabel{},
		&models.RepoTracking{},
		&models.LabelTracking{},
	)

	if err != nil {
		panic(err)
	}
}

//DbConnection returns the connection to use the db
func DbConnection() *gorm.DB {

	return db
}
