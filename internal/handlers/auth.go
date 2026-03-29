package handlers

import (
	"encoding/json"
	"net/http"
	"parus-test/internal/dto"
	"parus-test/internal/service"
)

type TokenResponse struct {
	Token string `json:"token"`
}

// SignIn godoc
// @Summary Авторизация
// @Description Возвращает токен
// @Tags auth
// @Produce json
// @Param request body dto.AuthRequest true "Логин и пароль"
// @Success 200  {object}  TokenResponse
// @Router /auth [post]
func SignIn(s service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var req dto.AuthRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		token, err := s.SignIn(r.Context(), req)
		if err != nil {
			respondError(w, err)
			return
		}
		resp := TokenResponse{
			Token: token,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
