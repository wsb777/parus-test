package domain

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type GroupRepo interface {
	Create(ctx context.Context, group *Group) (uuid.UUID, error)
	GetGroupList(ctx context.Context) ([]Group, error)
}

type UserRepo interface {
	Create(ctx context.Context, user User) error
	GetUserByName(ctx context.Context, username string) (User, error)
	GetUserList(ctx context.Context) ([]User, error)
}

type FileRepo interface {
	CreateEntry(ctx context.Context, fileID uuid.UUID, groupID uuid.UUID, name string, size int64) error
	FindEntry(ctx context.Context, fileID uuid.UUID) (*FileEntry, error)
	SetCurrentVersion(ctx context.Context, fileID uuid.UUID, semver SemVer) error

	SaveVersion(ctx context.Context, file *FileVersion) error
	FindVersion(ctx context.Context, fileID uuid.UUID, semver SemVer) (*FileVersion, error)
	ListVersions(ctx context.Context, fileID uuid.UUID, groupID uuid.UUID) ([]*FileVersion, error)
}

type FileStorage interface {
	Save(ctx context.Context, path string, data []byte) error
	Delete(ctx context.Context, path string) error
	Read(ctx context.Context, path string) (io.ReadCloser, error)
}

type TokenRepo interface {
	Create(ctx context.Context, token Token) error
	FindByHash(ctx context.Context, hash string) (*Token, error)
	Revoke(ctx context.Context, hash string) error
	RevokeAll(ctx context.Context, userID uuid.UUID) error
}
