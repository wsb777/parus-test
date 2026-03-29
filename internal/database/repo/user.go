package repo

import (
	"context"
	"errors"
	"fmt"
	"parus-test/internal/database/models"
	"parus-test/internal/domain"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

type UserRepo interface {
	Create(context.Context, domain.User) error
	GetUserByName(context.Context, string) (domain.User, error)
	GetUserList(context.Context) ([]domain.User, error)
}

func NewUserRepo(db *gorm.DB) domain.UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user domain.User) error {
	model := models.User{
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
		GroupID:  user.GroupID,
	}

	err := r.db.WithContext(ctx).Create(&model).Error

	if err != nil {
		return fmt.Errorf("user %q: %w", user.Username, mapDBError(err))
	}

	return err
}

func (r *userRepo) GetUserByName(ctx context.Context, username string) (domain.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("user %q: %w", username, domain.ErrUserNotFound)
		}

		return domain.User{}, err
	}

	return toDomainUser(user), nil
}

func (r *userRepo) GetUserList(ctx context.Context) ([]domain.User, error) {
	var users []models.User

	if err := r.db.WithContext(ctx).Select("id", "username", "group_id", "role").Preload("Group").Find(&users).Error; err != nil {
		return nil, err
	}

	var domainUsers []domain.User

	for _, user := range users {
		domainUsers = append(domainUsers, toDomainUser(user))
	}

	return domainUsers, nil
}

func toDomainUser(model models.User) domain.User {
	return domain.User{
		ID:       model.ID,
		Username: model.Username,
		Password: model.Password,
		Role:     model.Role,
		GroupID:  model.GroupID,
		Group: domain.Group{
			ID:   model.Group.ID,
			Name: model.Group.Name,
		},
	}
}
