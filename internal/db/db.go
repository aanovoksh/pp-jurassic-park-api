package db

import (
	"os"

	dbmodels "pp-jurassic-park-api/internal/db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsnMain = "host=postgres user=gorm password=gorm dbname=gorm sslmode=disable TimeZone=UTC"
const dsnTest = "host=postgrestest user=gormtest password=gormtest dbname=gormtest sslmode=disable TimeZone=UTC"

func Connect() (*gorm.DB, error) {
	if os.Getenv("GO_ENV") == "test" {
		return gorm.Open(postgres.Open(dsnTest), &gorm.Config{})
	}
	return gorm.Open(postgres.Open(dsnMain), &gorm.Config{})
}

func Migrate() error {
	dbConn, err := Connect()
	if err != nil {
		return err
	}
	return dbConn.AutoMigrate(&dbmodels.Cage{}, &dbmodels.Dinosaur{})
}
