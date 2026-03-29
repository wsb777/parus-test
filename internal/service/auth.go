package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"parus-test/internal/domain"
	"parus-test/internal/dto"
	"parus-test/pkg/hasher"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService interface {
	SignIn(ctx context.Context, req dto.AuthRequest) (string, error)
	CheckToken(ctx context.Context, raw string) (*domain.Token, error)
	RevokeAll(ctx context.Context, userID string) error
	RevokeToken(ctx context.Context, raw string) error
}

type authService struct {
	userRepo  domain.UserRepo
	tokenRepo domain.TokenRepo
	logger    *zap.Logger
}

func NewAuthService(userRepo domain.UserRepo, tokenRepo domain.TokenRepo, logger *zap.Logger) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}

func (s *authService) SignIn(ctx context.Context, req dto.AuthRequest) (string, error) {
	s.logger.Info("attempting to create sign in", zap.String("username", req.Username))

	user, err := s.userRepo.GetUserByName(ctx, req.Username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			s.logger.Warn("user not found", zap.String("username", req.Username))
		} else {
			s.logger.Error("failed to get user", zap.String("username", req.Username), zap.Error(err))
		}
		return "", fmt.Errorf("login or password are incorrect: %w", domain.ErrUnauthorized)
	}

	if !hasher.CheckPassword(req.Password, user.Password) {
		s.logger.Warn("invalid password", zap.String("username", req.Username))
		return "", fmt.Errorf("login or password are incorrect: %w", domain.ErrUnauthorized)
	}

	raw, hash, err := generateToken()
	if err != nil {
		s.logger.Error("failed to generate token", zap.String("username", req.Username), zap.Error(err))
		return "", fmt.Errorf("internal server error")
	}

	tokenDomain := domain.Token{
		ID:      uuid.New(),
		UserID:  user.ID,
		Role:    user.Role,
		GroupID: user.GroupID,
		Hash:    hash,
	}

	if err = s.tokenRepo.Create(ctx, tokenDomain); err != nil {
		s.logger.Error("failed to save token", zap.String("username", req.Username), zap.Error(err))
		return "", fmt.Errorf("internal server error")
	}

	return raw, nil
}

func (s *authService) CheckToken(ctx context.Context, raw string) (*domain.Token, error) {
	s.logger.Info("attemping to check token")
	hash := hashToken(raw)
	token, err := s.tokenRepo.FindByHash(ctx, hash)
	if err != nil {
		s.logger.Warn("token not found")
		return nil, err
	}
	return token, nil
}

func (s *authService) RevokeAll(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)

	if err != nil {
		s.logger.Warn("failed to parse uuid", zap.String("user_id", userID))
		return fmt.Errorf("internal server error")
	}
	err = s.tokenRepo.RevokeAll(ctx, id)

	if err != nil {
		s.logger.Warn("tokens not found")
		return fmt.Errorf("tokens not found: %w", domain.ErrTokenNotFound)
	}

	return nil
}

func (s *authService) RevokeToken(ctx context.Context, raw string) error {

	hash := hashToken(raw)
	err := s.tokenRepo.Revoke(ctx, hash)

	if err != nil {
		s.logger.Warn("tokens not found")
		return domain.ErrTokenNotFound
	}

	return nil
}

// Opaque token использую впервые

func generateToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	raw = hex.EncodeToString(b)
	hash = hashToken(raw)
	return
}

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
