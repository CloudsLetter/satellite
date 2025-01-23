package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"satellite/models/db"
	"sync"
	"time"
)

var (
	dbContext       *gorm.DB
	once            sync.Once
	autoMigrateOnce sync.Once
)

func InitDB(dsn string) {
	once.Do(func() {
		dbc, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		sqlDB, err := dbc.DB()
		sqlDB.SetMaxIdleConns(2)
		sqlDB.SetMaxOpenConns(4)
		sqlDB.SetConnMaxLifetime(time.Hour)
		dbContext = dbc
	})
}

func AutoMigrate(dbc *gorm.DB) {
	autoMigrateOnce.Do(func() {
		err := dbc.AutoMigrate(&db.CallsStatistics{})
		if err != nil {
			panic(err)
		}
	})
}
func GetDBContext() *gorm.DB {
	return dbContext
}
