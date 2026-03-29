package dto

import (
	"io"
	"parus-test/internal/domain"
	"time"
)

type UserRequest struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Group    string `json:"group"`
	Role     string `json:"role"`
}

type GroupRequest struct {
	Name string `json:"name"`
}

type GroupResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FileUploadRequest struct {
	ID          string `json:"id"`
	Name        string `json:"filename"`
	Content     []byte `json:"content"`
	ContentType string `json:"content_type"`
	GroupID     string
	Version     string          `json:"version"`
	Bump        domain.BumpType `json:"bump"`
}

type FileVersionResponse struct {
	FileID    string    `json:"file_id"`
	Version   string    `json:"version"`
	Checksum  string    `json:"checksum,omitempty"`
	Size      int64     `json:"size"`
	URL       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type FileInfoRequest struct {
	FileID  string
	GroupID string
}

type FileInfoResponse struct {
	FileID   string   `json:"file_id"`
	Versions []string `json:"versions"`
}

type FileDataRequest struct {
	FileID      string
	FileVersion string
	GroupID     string
}

type FileDataResponse struct {
	Hash        string
	Reader      io.ReadCloser
	ContentType string
	FileName    string
	Size        int64
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserListItemResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
	Role      string `json:"role"`
}
