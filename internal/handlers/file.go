package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"parus-test/internal/domain"
	"parus-test/internal/dto"
	"parus-test/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// UploadFile godoc
// @Summary Загрузка файла
// @Description Возвращает информацию о загруженном файле
// @Tags file
// @Param file_id path string false "ID файла"
// @Produce json
// @Success 200 {object} dto.FileVersionResponse
// @Router /admin/files/{file_id} [post]
func UploadFile(s service.FileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		r.ParseMultipartForm(10 << 20)

		groupID := r.FormValue("group_id")
		if groupID == "" {
			http.Error(w, "group_id required", http.StatusBadRequest)
			return
		}

		bump := r.FormValue("bump")
		if bump == "" {
			http.Error(w, "bump required", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "failed to read file", http.StatusInternalServerError)
			return
		}

		fileDTO := dto.FileUploadRequest{
			Content:     data,
			Name:        header.Filename,
			ContentType: header.Header.Get("Content-Type"),
			GroupID:     groupID,
			ID:          chi.URLParam(r, "file_id"),
			Bump:        domain.BumpType(bump),
		}

		result, err := s.UploadNewVersion(r.Context(), fileDTO)
		if err != nil {
			respondError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

// GetFileInfo godoc
// @Summary Получение информации о версиях файла
// @Description Возвращает информацию о версиях файла
// @Tags file
// @Param file_id path string true "ID файла"
// @Produce json
// @Success 200 {object} dto.FileVersionResponse
// @Router /files/{file_id}/version/list [get]
func GetFileInfo(s service.FileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.FileInfoRequest
		req.FileID = chi.URLParam(r, "file_id")
		if req.FileID == "" {
			http.Error(w, "file_id required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		groupID, ok := ctx.Value("group_id").(string)
		if !ok || groupID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		req.GroupID = groupID

		result, err := s.GetFileInfo(ctx, req)
		if err != nil {
			respondError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

// GetFileLatestVersionInfo godoc
// @Summary Получение информации о файле последнего обновления
// @Description Возвращает информацию о последней версии
// @Param file_id path string true "ID файла"
// @Tags file
// @Produce json
// @Success 200 {object} dto.FileVersionResponse
// @Router /files/{file_id}/version/latest [get]
func GetFileLatestVersionInfo(s service.FileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.FileInfoRequest
		req.FileID = chi.URLParam(r, "file_id")
		if req.FileID == "" {
			http.Error(w, "file_id required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		groupID, ok := ctx.Value("group_id").(string)
		if !ok || groupID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		req.GroupID = groupID

		result, err := s.GetFileLatestVersionInfo(ctx, req)
		if err != nil {
			respondError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

// GetFileData godoc
// @Summary Получение файла
// @Description Возращает файл и хеш в header Digest
// @Tags file
// @Param file_id path string true "ID файла"
// @Param file_version path string true "Версия файла"
// @Produce octet-stream
// @Success 200 {file} binary
// @Router /files/{file_id}/version/{file_version}/download [get]
func GetFileData(s service.FileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.FileDataRequest
		req.FileID = chi.URLParam(r, "file_id")
		if req.FileID == "" {
			http.Error(w, "file_id required", http.StatusBadRequest)
			return
		}

		req.FileVersion = chi.URLParam(r, "file_version")
		if req.FileVersion == "" {
			http.Error(w, "file_version required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		groupID, ok := ctx.Value("group_id").(string)
		if !ok || groupID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		req.GroupID = groupID

		result, err := s.GetFileDataByVersion(ctx, req)
		if err != nil {
			respondError(w, err)
			return
		}

		w.Header().Set("Content-Type", result.ContentType)
		w.Header().Set("Content-Disposition", "attachment; filename="+result.FileName)
		w.Header().Set("Content-Length", strconv.FormatInt(result.Size, 10))
		w.Header().Set("Digest", result.Hash)
		io.Copy(w, result.Reader)
	}
}
