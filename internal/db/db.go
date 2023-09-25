package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=postgres user=gorm password=gorm dbname=gorm sslmode=disable TimeZone=UTC"

func Connect() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
