package handlers

import (
	"encoding/json"
	"net/http"
	"parus-test/internal/dto"
	"parus-test/internal/service"

	"github.com/go-chi/chi/v5"
)

// CreateUser godoc
// @Summary Создание пользователя
// @Description Возращает статус в header
// @Tags users
// @Produce json
// @Param request body dto.UserRequest true "Логин, пароль и группа"
// @Success 201
// @Router /admin/users [post]
func CreateUser(s service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.UserRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if err := s.CreateUser(r.Context(), req); err != nil {
			respondError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

// GetUserList godoc
// @Summary Получение списка пользователей
// @Description Возвращает список пользователей с их группами
// @Tags users
// @Produce json
// @Success 200 {array} dto.UserListItemResponse
// @Router /admin/users [get]
func GetUserList(s service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		users, err := s.GetUserList(r.Context())
		if err != nil {
			respondError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// RevokeUserTokenAll godoc
// @Summary Отзыв всех токенов пользователя
// @Description Возращает статус в header
// @Tags users
// @Param user_id path string true "ID пользователя"
// @Success 200
// @Router /admin/users/{user_id}/token/revoke/all [post]
func RevokeUserTokenAll(s service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "user_id")
		if userID == "" {
			http.Error(w, "user_id required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		err := s.RevokeAll(ctx, userID)
		if err != nil {
			respondError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// RevokeUserToken godoc
// @Summary Отзыв конкретного токена
// @Description Возращает статус в header
// @Tags users
// @Param token_raw path string true "Токен"
// @Success 200
// @Router /admin/users/token/{token_raw}/revoke [post]
func RevokeUserToken(s service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		raw := chi.URLParam(r, "token_raw")
		if raw == "" {
			http.Error(w, "token_raw required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		err := s.RevokeToken(ctx, raw)
		if err != nil {
			respondError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
