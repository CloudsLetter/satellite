package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"satellite/models/db"
	"time"
)

var dbContext *gorm.DB

func InitDB(dsn string) {
	dbc, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := dbc.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(4)
	sqlDB.SetConnMaxLifetime(time.Hour)
	dbContext = dbc

}

func AutoMigrate(dbc *gorm.DB) {
	err := dbc.AutoMigrate(&db.CallsStatistics{})
	if err != nil {
		panic(err)
	}
}
func GetDBContext() *gorm.DB {
	return dbContext
}
