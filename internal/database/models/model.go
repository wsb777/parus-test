package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

type User struct {
	BaseModel
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	GroupID   uuid.UUID `gorm:"not null"`
	Role      string    `gorm:"not null"`
	CreatedAt time.Time

	Group Group `gorm:"foreignKey:GroupID"`
}

type Group struct {
	BaseModel
	Name      string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time

	Users []User      `gorm:"foreignKey:GroupID"`
	Files []FileEntry `gorm:"foreignKey:GroupID"`
}

type FileEntry struct {
	BaseModel
	Name           string `gorm:"unique;not null"`
	CurrentVersion string `gorm:"not null"`
	CreatedAt      time.Time
	GroupID        uuid.UUID `gorm:"not null"`
	Size           int64     `gorm:"not null"`

	Group    Group         `gorm:"foreignKey:GroupID"`
	Versions []FileVersion `gorm:"foreignKey:FileID"`
}

type FileVersion struct {
	BaseModel
	FileID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_file_version"`
	Version     string    `gorm:"not null;uniqueIndex:idx_file_version"`
	StoragePath string    `gorm:"not null"`
	Checksum    string    `gorm:"not null"`
	Size        int64     `gorm:"not null"`
	CreatedAt   time.Time
	FileEntry   FileEntry `gorm:"foreignKey:FileID"`
}

type Token struct {
	BaseModel
	Hash       string    `gorm:"unique"`
	UserID     uuid.UUID `gorm:"not null"`
	CreatedAt  time.Time
	LastUsedAt *time.Time
	RevokedAt  *time.Time

	User User `gorm:"foreignKey:UserID"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}
