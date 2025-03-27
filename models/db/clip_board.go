package db

import (
	"gorm.io/gorm"
	"time"
)

type ClipBoard struct {
	gorm.Model
	ID                 int64 `gorm:"primaryKey,autoIncrement"`
	Title              string
	Content            string
	SyntaxHighlighting string
	Expire             time.Time
	Snapchat           bool
}
