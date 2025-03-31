package db

import (
	"gorm.io/gorm"
	"time"
)

type ClipBoard struct {
	gorm.Model         `json:"-"`
	ID                 int64     `gorm:"primaryKey,autoIncrement" json:"-"`
	Title              string    `json:"title"`
	Content            string    `json:"content"`
	SyntaxHighlighting string    `json:"syntax_highlighting"`
	Expire             time.Time `json:"expire"`
	Read               int       `json:"read"`
	Snapchat           bool      `json:"snapchat"`
	Password           string    `json:"-"`
}
