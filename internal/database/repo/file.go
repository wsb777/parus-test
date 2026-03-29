package repo

import (
	"context"
	"errors"
	"fmt"
	"parus-test/internal/database/models"
	"parus-test/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fileRepo struct {
	db *gorm.DB
}

func NewFileRepo(db *gorm.DB) domain.FileRepo {
	return &fileRepo{db: db}
}

func (r *fileRepo) CreateEntry(ctx context.Context, fileID uuid.UUID, groupID uuid.UUID, name string, size int64) error {
	model := models.FileEntry{
		BaseModel:      models.BaseModel{ID: fileID},
		GroupID:        groupID,
		Name:           name,
		CurrentVersion: "0.0.0",
		Size:           size,
	}

	result := r.db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		return fmt.Errorf("create file entry: %w", result.Error)
	}

	return nil
}

func (r *fileRepo) FindEntry(ctx context.Context, fileID uuid.UUID) (*domain.FileEntry, error) {
	var model models.FileEntry

	result := r.db.WithContext(ctx).Where("id = ?", fileID).First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("file entry %s: %w", fileID, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("find file entry: %w", result.Error)
	}

	return fileEntryToDomain(&model), nil
}

func (r *fileRepo) SaveVersion(ctx context.Context, file *domain.FileVersion) error {

	model := models.FileVersion{
		BaseModel:   models.BaseModel{ID: uuid.MustParse(file.ID)},
		FileID:      uuid.MustParse(file.FileID),
		Version:     file.Version.String(),
		StoragePath: file.StoragePath,
		Checksum:    file.Checksum,
		Size:        file.Size,
		CreatedAt:   file.CreatedAt,
	}

	result := r.db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		return fmt.Errorf("save file version: %w", result.Error)
	}

	return nil
}

func (r *fileRepo) FindVersion(ctx context.Context, fileID uuid.UUID, semver domain.SemVer) (*domain.FileVersion, error) {
	var model models.FileVersion

	result := r.db.WithContext(ctx).Where("file_id = ? AND version = ?", fileID, semver.String()).First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("version %s for file %s: %w", semver.String(), fileID, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("find file version: %w", result.Error)
	}
	return fileVersionToDomain(&model), nil
}

func (r *fileRepo) ListVersions(ctx context.Context, fileID uuid.UUID, groupID uuid.UUID) ([]*domain.FileVersion, error) {
	var models []models.FileVersion

	q := r.db.WithContext(ctx).Preload("FileEntry").
		Where("file_id = ?", fileID).
		Order("created_at ASC")

	if groupID != uuid.Nil {
		q = q.Joins("JOIN file_entries ON file_entries.id = file_versions.file_id").
			Where("file_entries.group_id = ?", groupID)
	}

	if err := q.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("list file versions: %w", err)
	}

	versions := make([]*domain.FileVersion, len(models))
	for i, m := range models {
		versions[i] = fileVersionToDomain(&m)
	}

	return versions, nil
}

func (r *fileRepo) SetCurrentVersion(ctx context.Context, fileID uuid.UUID, semver domain.SemVer) error {
	result := r.db.WithContext(ctx).
		Model(&models.FileEntry{}).
		Where("id = ?", fileID).
		Update("current_version", semver.String())
	if result.Error != nil {
		return fmt.Errorf("set current version: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("file entry %s: %w", fileID, domain.ErrNotFound)
	}
	return nil
}

func fileEntryToDomain(m *models.FileEntry) *domain.FileEntry {
	semver, err := domain.ParseSemVer(m.CurrentVersion)
	if err != nil {
		semver = domain.SemVer{}
	}
	return &domain.FileEntry{
		ID:             m.ID.String(),
		GroupID:        m.GroupID.String(),
		Name:           m.Name,
		CurrentVersion: semver,
		Size:           m.Size,
		CreatedAt:      m.CreatedAt,
	}
}

func fileVersionToDomain(model *models.FileVersion) *domain.FileVersion {
	semver, err := domain.ParseSemVer(model.Version)
	if err != nil {
		semver = domain.SemVer{}
	}
	return &domain.FileVersion{
		ID:          model.ID.String(),
		FileID:      model.FileID.String(),
		Version:     semver,
		StoragePath: model.StoragePath,
		Checksum:    model.Checksum,
		Size:        model.Size,
		CreatedAt:   model.CreatedAt,
	}
}
