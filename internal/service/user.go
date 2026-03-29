package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"parus-test/internal/domain"
	"parus-test/internal/dto"
	"parus-test/pkg/hasher"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userService struct {
	repo   domain.UserRepo
	logger *zap.Logger
}

type UserService interface {
	CreateUser(context.Context, dto.UserRequest) error
	CreateAdmin(string, string, uuid.UUID) error
	GetUserList(ctx context.Context) ([]dto.UserListItemResponse, error)
}

func NewUserService(repo domain.UserRepo, logger *zap.Logger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) CreateUser(ctx context.Context, req dto.UserRequest) error {
	s.logger.Info("attempting to create user", zap.String("username", req.Username))

	if strings.TrimSpace(req.Username) == "" {
		return fmt.Errorf("username cannot be empty: %w", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("password cannot be empty: %w", domain.ErrInvalidInput)
	}

	groupID, err := uuid.Parse(req.Group)
	if err != nil {
		s.logger.Warn("invalid group id", zap.String("group_id", req.Group))
		return fmt.Errorf("invalid group id: %w", domain.ErrInvalidInput)
	}

	hash, err := hasher.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", zap.String("username", req.Username), zap.Error(err))
		return fmt.Errorf("internal server error")
	}

	user := domain.User{
		Username: req.Username,
		Password: hash,
		Role:     req.Role,
		GroupID:  groupID,
	}

	if err = s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			s.logger.Warn("user already exists", zap.String("username", req.Username))
			return err
		}
		s.logger.Error("failed to create user", zap.String("username", req.Username), zap.Error(err))
		return fmt.Errorf("internal server error")
	}
	return nil
}

func (s *userService) CreateAdmin(username string, password string, groupID uuid.UUID) error {
	s.logger.Info("attempting to create admin", zap.String("username", username))

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	hash, err := hasher.HashPassword(password)
	if err != nil {
		s.logger.Error("failed to hash password", zap.String("username", username), zap.Error(err))
		return fmt.Errorf("internal server error")
	}

	user := domain.User{
		Username: username,
		Password: hash,
		GroupID:  groupID,
		Role:     "admin",
	}

	if err = s.repo.Create(ctx, user); err != nil {

		if errors.Is(err, domain.ErrAlreadyExists) {
			s.logger.Warn("admin already exist", zap.String("username", username))
			return err
		}

		s.logger.Error("failed to create admin", zap.String("username", username), zap.Error(err))
		return fmt.Errorf("internal server error")
	}
	return nil
}

func (s *userService) GetUserList(ctx context.Context) ([]dto.UserListItemResponse, error) {
	userModels, err := s.repo.GetUserList(ctx)

	if err != nil {
		s.logger.Error("failed to list users", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	users := make([]dto.UserListItemResponse, 0, len(userModels))

	for _, user := range userModels {
		users = append(users, dto.UserListItemResponse{
			ID:        user.ID.String(),
			Username:  user.Username,
			GroupID:   user.Group.ID.String(),
			GroupName: user.Group.Name,
			Role:      user.Role,
		})
	}

	return users, nil
}
