package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"parus-test/internal/domain"
	"parus-test/internal/dto"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type groupService struct {
	repo   domain.GroupRepo
	logger *zap.Logger
}

type GroupService interface {
	Create(context.Context, dto.GroupRequest) (string, error)
	CreateAdminGroup() (uuid.UUID, error)
	GetGroupList(context.Context) ([]dto.GroupResponse, error)
}

func NewGroupService(repo domain.GroupRepo, logger *zap.Logger) GroupService {
	return &groupService{
		repo:   repo,
		logger: logger,
	}
}

func (s *groupService) Create(ctx context.Context, groupDTO dto.GroupRequest) (string, error) {
	s.logger.Info("attempting to create group", zap.String("name", groupDTO.Name))

	if strings.TrimSpace(groupDTO.Name) == "" {
		return "", fmt.Errorf("group name cannot be empty: %w", domain.ErrInvalidInput)
	}

	group := &domain.Group{Name: groupDTO.Name}
	id, err := s.repo.Create(ctx, group)

	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			s.logger.Warn("group already exists", zap.String("name", groupDTO.Name))
			return "", err
		}
		s.logger.Error("failed to create group", zap.String("name", groupDTO.Name), zap.Error(err))
		return "", fmt.Errorf("internal server error")
	}

	return id.String(), nil
}

func (s *groupService) CreateAdminGroup() (uuid.UUID, error) {
	s.logger.Info("attempting to create admin group")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	group := &domain.Group{Name: "admin"}
	id, err := s.repo.Create(ctx, group)

	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			s.logger.Warn("admin group already exists")
			return uuid.Nil, err
		}
		s.logger.Error("failed to create admin group", zap.Error(err))
		return uuid.Nil, fmt.Errorf("internal server error")
	}

	return id, nil
}

func (s *groupService) GetGroupList(ctx context.Context) ([]dto.GroupResponse, error) {
	s.logger.Info("attempting to get group list")

	groupModels, err := s.repo.GetGroupList(ctx)

	if err != nil {
		s.logger.Error("failed to list groups", zap.Error(err))
		return nil, fmt.Errorf("internal server error")
	}

	groups := make([]dto.GroupResponse, 0, len(groupModels))

	for _, groupModel := range groupModels {
		groups = append(groups, dto.GroupResponse{
			ID:   groupModel.ID.String(),
			Name: groupModel.Name,
		})
	}

	return groups, nil
}
