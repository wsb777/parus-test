package repo

import (
	"context"
	"errors"
	"parus-test/internal/database/models"
	"parus-test/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) domain.TokenRepo {
	return &tokenRepo{
		db: db,
	}
}

func (r *tokenRepo) Create(ctx context.Context, token domain.Token) error {
	m := models.Token{
		BaseModel: models.BaseModel{ID: token.ID},
		Hash:      token.Hash,
		UserID:    token.UserID,
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *tokenRepo) FindByHash(ctx context.Context, hash string) (*domain.Token, error) {
	var m models.Token

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("hash = ? AND revoked_at IS NULL", hash).
		First(&m).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	r.db.WithContext(ctx).Model(&m).Update("last_used_at", gorm.Expr("now()"))

	return toDomain(m), nil
}

func (r *tokenRepo) Revoke(ctx context.Context, hash string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Token{}).
		Where("hash = ? AND revoked_at IS NULL", hash).
		Update("revoked_at", gorm.Expr("now()"))

	if result.RowsAffected == 0 {
		return domain.ErrTokenNotFound
	}
	return result.Error
}

func (r *tokenRepo) RevokeAll(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Token{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", gorm.Expr("now()")).Error
}

func toDomain(m models.Token) *domain.Token {
	return &domain.Token{
		ID:         m.ID,
		Hash:       m.Hash,
		UserID:     m.UserID,
		Role:       m.User.Role,
		GroupID:    m.User.GroupID,
		CreatedAt:  m.CreatedAt,
		LastUsedAt: m.LastUsedAt,
		RevokedAt:  m.RevokedAt,
	}
}
