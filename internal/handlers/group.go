package handlers

import (
	"encoding/json"
	"net/http"
	"parus-test/internal/dto"
	"parus-test/internal/service"
)

// CreateGroup godoc
// @Summary Создание группы
// @Description Возвращает id группы
// @Tags group
// @Produce json
// @Param request body dto.GroupRequest true "Наименование"
// @Success 201 {object} dto.GroupResponse
// @Router /admin/groups [post]
func CreateGroup(s service.GroupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var req dto.GroupRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		id, err := s.Create(r.Context(), req)
		if err != nil {
			respondError(w, err)
			return
		}

		resp := dto.GroupResponse{
			Name: req.Name,
			ID:   id,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

// GetGroupList godoc
// @Summary Получение списка групп
// @Description Возвращает список групп
// @Tags group
// @Produce json
// @Success 200 array  dto.GroupResponse
// @Router /admin/groups [get]
func GetGroupList(s service.GroupService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		groups, err := s.GetGroupList(r.Context())
		if err != nil {
			respondError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}
