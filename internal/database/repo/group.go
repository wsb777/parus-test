package repo

import (
	"context"
	"fmt"
	"parus-test/internal/database/models"
	"parus-test/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type groupRepo struct {
	db *gorm.DB
}

func NewGroupRepo(db *gorm.DB) domain.GroupRepo {
	return &groupRepo{
		db: db,
	}
}

// TODO
// Gorm не логирует ошибки postgres о дублирующихся ключах, из за этого логи засоряются
func (r *groupRepo) Create(ctx context.Context, group *domain.Group) (uuid.UUID, error) {
	model := models.Group{
		Name: group.Name,
	}
	err := r.db.WithContext(ctx).Create(&model).Error

	if err != nil {
		return uuid.Nil, fmt.Errorf("group %q: %w", group.Name, mapDBError(err))
	}

	return model.ID, err
}

func (r *groupRepo) GetGroupList(ctx context.Context) ([]domain.Group, error) {
	var groups []models.Group

	if err := r.db.WithContext(ctx).Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("groups: %w", mapDBError(err))
	}

	var domainGroups []domain.Group

	for _, group := range groups {
		domainGroups = append(domainGroups, toDomainGroup(group))
	}

	return domainGroups, nil
}

func toDomainGroup(model models.Group) domain.Group {
	return domain.Group{
		ID:   model.ID,
		Name: model.Name,
	}
}
