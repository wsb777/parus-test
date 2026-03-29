package domain

import (
	"time"

	"github.com/google/uuid"
)

//  В рамках тестового задания не делал отдельную таблицу под роли

type User struct {
	ID       uuid.UUID
	Username string
	Password string
	GroupID  uuid.UUID
	Role     string

	Group Group
}

type Group struct {
	ID   uuid.UUID
	Name string
}

type FileEntry struct {
	ID             string
	Name           string
	GroupID        string
	Size           int64
	CurrentVersion SemVer
	CreatedAt      time.Time
}

// Возможно стоило бы добавить owner файла
type FileVersion struct {
	ID          string
	FileID      string
	Version     SemVer
	StoragePath string
	Checksum    string
	Size        int64
	CreatedAt   time.Time
}

type Token struct {
	ID         uuid.UUID
	Hash       string
	UserID     uuid.UUID
	GroupID    uuid.UUID
	Role       string
	CreatedAt  time.Time
	LastUsedAt *time.Time
	RevokedAt  *time.Time
}
