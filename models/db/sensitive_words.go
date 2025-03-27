package db

import "gorm.io/gorm"

type SensitiveWords struct {
	gorm.Model
	ID       int64  `gorm:"primaryKey,autoIncrement"`
	Keywords string `gorm:"not null;index:Keywords,unique"`
}
