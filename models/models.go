package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username string    `gorm:"uniqueIndex;not null"`
	Password string    `gorm:"not null"`
	Contents []Content `gorm:"foreignKey:UserId"`
}

type Content struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Link   string    `gorm:"not null"`
	Title  string    `gorm:"not null"`
	Type   string    `gorm:"not null"`
	UserId uuid.UUID `gorm:"type:uuid; not null"`

	Tags []Tag `gorm:"many2many:content_tags;"`
}

type Tag struct {
	ID    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title string    `gorm:"not null"`
}
