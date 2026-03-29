package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BumpType string

const (
	BumpPatch BumpType = "patch"
	BumpMinor BumpType = "minor"
	BumpMajor BumpType = "major"
)

func NewFileVersion(groupID string, fileID string, current SemVer, bump BumpType, data []byte, ext string) *FileVersion {

	next := bumpVersion(current, bump)
	id := uuid.New().String()
	path := fmt.Sprintf("%s/%s/%s/%s%s", groupID, fileID, next.String(), id, ext)

	return &FileVersion{
		ID:          id,
		FileID:      fileID,
		Version:     next,
		StoragePath: path,
		Checksum:    computeSHA256(data),
		Size:        int64(len(data)),
		CreatedAt:   time.Now(),
	}
}

func computeSHA256(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func bumpVersion(v SemVer, bump BumpType) SemVer {
	switch bump {
	case BumpMinor:
		return v.BumpMinor()
	case BumpMajor:
		return v.BumpMajor()
	default:
		return v.BumpPatch()
	}
}
