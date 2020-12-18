module github.com/nocubicles/develytica

go 1.15

require (
	github.com/aws/aws-sdk-go v1.35.33
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/google/go-github/v33 v33.0.0
	github.com/gorilla/mux v1.8.0
	github.com/jinzhu/gorm v1.9.16
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.8.0 // indirect
	golang.org/x/oauth2 v0.0.0-20201203001011-0b49973bad19
	gorm.io/driver/postgres v1.0.5
	gorm.io/gorm v1.20.8
)

replace github.com/google/go-github/v33 => github.com/nocubicles/go-github/v33 v33.0.1-0.20201216161729-92013ff0a7d4
