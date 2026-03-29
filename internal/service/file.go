package service

import (
	"context"
	"errors"
	"fmt"
	"parus-test/internal/domain"
	"parus-test/internal/dto"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nlgolib/mime"
	"go.uber.org/zap"
)

type fileService struct {
	repo    domain.FileRepo
	storage domain.FileStorage
	logger  *zap.Logger
}

type FileService interface {
	UploadNewVersion(ctx context.Context, file dto.FileUploadRequest) (*dto.FileVersionResponse, error)
	GetFileInfo(ctx context.Context, file dto.FileInfoRequest) (*dto.FileInfoResponse, error)
	GetFileLatestVersionInfo(ctx context.Context, file dto.FileInfoRequest) (*dto.FileVersionResponse, error)
	GetFileDataByVersion(ctx context.Context, file dto.FileDataRequest) (*dto.FileDataResponse, error)
}

func NewFileService(repo domain.FileRepo, storage domain.FileStorage, logger *zap.Logger) FileService {
	return &fileService{
		repo:    repo,
		storage: storage,
		logger:  logger,
	}
}

func (s *fileService) UploadNewVersion(ctx context.Context, file dto.FileUploadRequest) (*dto.FileVersionResponse, error) {
	s.logger.Info("attempting to uoload file version", zap.String("file_id", file.ID), zap.String("group_id", file.GroupID))

	groupUUID, err := uuid.Parse(file.GroupID)
	if err != nil {
		s.logger.Warn("invalid group id", zap.String("group_id", file.GroupID))
		return nil, domain.ErrInvalidInput
	}

	var currentVersion domain.SemVer

	if file.ID == "" {
		newID := uuid.New()
		file.ID = newID.String()

		if err := s.repo.CreateEntry(ctx, newID, groupUUID, file.Name, int64(len(file.Content))); err != nil {
			s.logger.Error("failed to create file entry", zap.String("file_id", file.ID), zap.Error(err))
			return nil, errors.New("internal server error")
		}
	} else {
		id, err := uuid.Parse(file.ID)

		if err != nil {
			s.logger.Warn("invalid file id", zap.String("file_id", file.ID))
			return nil, domain.ErrInvalidInput
		}

		entry, err := s.repo.FindEntry(ctx, id)

		if err != nil {
			s.logger.Error("failed to find file entry", zap.String("file_id", file.ID), zap.Error(err))
			return nil, domain.ErrNotFound
		}

		if entry.GroupID != file.GroupID {
			s.logger.Warn("access denied",
				zap.String("file_id", file.ID),
				zap.String("file_group", entry.GroupID),
				zap.String("caller_group", file.GroupID),
			)
			return nil, domain.ErrForbidden
		}

		currentVersion = entry.CurrentVersion
	}

	ext := mime.MimeToExt(file.ContentType)
	ver := domain.NewFileVersion(file.GroupID, file.ID, currentVersion, file.Bump, file.Content, ext)

	if err := s.storage.Save(ctx, ver.StoragePath, file.Content); err != nil {
		s.logger.Error("failed to save file to storage", zap.String("file_id", file.ID), zap.Error(err))
		return nil, errors.New("internal server error")
	}

	if err := s.repo.SaveVersion(ctx, ver); err != nil {
		s.storage.Delete(ctx, ver.StoragePath)
		s.logger.Error("failed to save file version", zap.String("file_id", file.ID), zap.Error(err))
		return nil, errors.New("internal server error")
	}

	fileID, err := uuid.Parse(ver.FileID)

	if err != nil {
		s.logger.Error("failed to parse file_id", zap.String("file_id", ver.FileID), zap.Error(err))
		return nil, errors.New("internal server error")
	}

	if err := s.repo.SetCurrentVersion(ctx, fileID, ver.Version); err != nil {
		s.logger.Error("failed to set current version", zap.String("file_id", file.ID), zap.Error(err))
		return nil, errors.New("internal server error")
	}

	return &dto.FileVersionResponse{
		FileID:    ver.FileID,
		Version:   ver.Version.String(),
		Checksum:  ver.Checksum,
		Size:      ver.Size,
		CreatedAt: ver.CreatedAt,
	}, nil
}

func (s *fileService) GetFileInfo(ctx context.Context, file dto.FileInfoRequest) (*dto.FileInfoResponse, error) {
	s.logger.Info("attempting to get file info", zap.String("file_id", file.FileID), zap.String("group_id", file.GroupID))

	fileID, err := uuid.Parse(file.FileID)

	if err != nil {
		s.logger.Warn("invalid file id", zap.String("file_id", file.FileID))
		return nil, domain.ErrInvalidInput
	}

	entry, err := s.repo.FindEntry(ctx, fileID)

	if err != nil {
		s.logger.Error("failed to find file entry", zap.String("file_id", file.FileID), zap.Error(err))
		return nil, fmt.Errorf("find entry: %w", err)
	}

	if entry.GroupID != file.GroupID {
		s.logger.Warn("access denied",
			zap.String("file_id", file.FileID),
			zap.String("file_group", entry.GroupID),
			zap.String("caller_group", file.GroupID),
		)
		return nil, domain.ErrForbidden
	}

	listVersions, err := s.repo.ListVersions(ctx, fileID, uuid.Nil)
	if err != nil {
		s.logger.Error("failed to list file versions", zap.String("file_id", file.FileID), zap.Error(err))
		return nil, errors.New("internal server error")
	}

	versions := make([]string, 0, len(listVersions))
	for _, ver := range listVersions {
		versions = append(versions, ver.Version.String())
	}

	return &dto.FileInfoResponse{
		FileID:   file.FileID,
		Versions: versions,
	}, nil
}

func (s *fileService) GetFileLatestVersionInfo(ctx context.Context, file dto.FileInfoRequest) (*dto.FileVersionResponse, error) {
	s.logger.Info("attempting to get file lasted info", zap.String("file_id", file.FileID), zap.String("group_id", file.GroupID))

	fileID, err := uuid.Parse(file.FileID)

	if err != nil {
		s.logger.Warn("invalid file id", zap.String("file_id", file.FileID))
		return nil, domain.ErrInvalidInput
	}

	entry, err := s.repo.FindEntry(ctx, fileID)

	if err != nil {
		s.logger.Error("failed to find file entry", zap.String("file_id", file.FileID), zap.Error(err))
		return nil, fmt.Errorf("find entry: %w", err)
	}

	if entry.GroupID != file.GroupID {
		s.logger.Warn("access denied",
			zap.String("file_id", file.FileID),
			zap.String("file_group", entry.GroupID),
			zap.String("caller_group", file.GroupID),
		)
		return nil, domain.ErrForbidden
	}

	return &dto.FileVersionResponse{
		FileID:    entry.ID,
		Version:   entry.CurrentVersion.String(),
		Size:      entry.Size,
		CreatedAt: entry.CreatedAt,
	}, nil
}

func (s *fileService) GetFileDataByVersion(ctx context.Context, file dto.FileDataRequest) (*dto.FileDataResponse, error) {
	s.logger.Info("attempting to get file data by version", zap.String("file_id", file.FileID), zap.String("group_id", file.GroupID), zap.String("version", file.FileVersion))

	fileID, err := uuid.Parse(file.FileID)
	if err != nil {
		s.logger.Warn("invalid file id", zap.String("file_id", file.FileID))
		return nil, domain.ErrInvalidInput
	}

	semver, err := domain.ParseSemVer(file.FileVersion)
	if err != nil {
		s.logger.Warn("invalid  id", zap.String("file_version", file.FileVersion))
		return nil, domain.ErrInvalidInput
	}

	entry, err := s.repo.FindEntry(ctx, fileID)

	if err != nil {
		s.logger.Error("failed to find file entry", zap.String("file_id", file.FileID), zap.Error(err))
		return nil, domain.ErrNotFound
	}

	if entry.GroupID != file.GroupID {
		s.logger.Warn("access denied",
			zap.String("file_id", file.FileID),
			zap.String("file_group", entry.GroupID),
			zap.String("caller_group", file.GroupID),
		)
		return nil, domain.ErrForbidden
	}

	ver, err := s.repo.FindVersion(ctx, fileID, semver)
	if err != nil {
		s.logger.Error("failed to find file version", zap.String("file_id", file.FileID), zap.String("file_version", file.FileVersion), zap.Error(err))
		return nil, domain.ErrNotFound
	}

	reader, err := s.storage.Read(ctx, ver.StoragePath)

	if err != nil {
		s.logger.Error("failed to find file version", zap.String("file_id", file.FileID), zap.String("file_version", file.FileVersion), zap.Error(err))
		return nil, errors.New("internal server error")
	}

	contentType := mime.ExtToMime(filepath.Ext(ver.StoragePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return &dto.FileDataResponse{
		Hash:        ver.Checksum,
		ContentType: contentType,
		Reader:      reader,
		FileName:    entry.Name + "-v" + semver.String(),
		Size:        ver.Size,
	}, nil
}
