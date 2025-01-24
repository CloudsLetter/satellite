package db

import (
	"gorm.io/gorm"
)

type CallsStatistics struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey"`
	Minute      int    `gorm:"not null;index:minute_index"`
	Hour        int    `gorm:"not null;index:hour_index"`
	Day         int    `gorm:"not null;index:day_index"`
	Month       int    `gorm:"not null;index:month_index"`
	Year        int    `gorm:"not null;index:year_index"`
	Source      string `gorm:"not null;index:source_index"`
	Destination string `gorm:"not null;index:destination_index"`
	UserAgent   string `gorm:"not null"`
	Country     string `gorm:"not null"`
	Region      string `gorm:"not null"`
	Province    string `gorm:"not null"`
	City        string `gorm:"not null"`
	ISP         string `gorm:"not null"`
}
