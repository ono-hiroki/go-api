// Package user はユーザー関連のHTTPハンドラーを提供する。
package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-api/internal/application/user"
	httperrors "go-api/internal/presentation/http/errors"
)

// listUserResponse はユーザー情報のJSONレスポンス。
type listUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// listUsersResponse はユーザー一覧のJSONレスポンス。
type listUsersResponse struct {
	Users []listUserResponse `json:"users"`
}

func newListUsersResponse(output *user.ListUsersOutput) listUsersResponse {
	users := make([]listUserResponse, len(output.Users))
	for i, u := range output.Users {
		users[i] = listUserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}
	}
	return listUsersResponse{Users: users}
}

// ListHandler はユーザー一覧取得のHTTPハンドラー。
type ListHandler struct {
	uc     *user.ListUsersUsecase
	logger *slog.Logger
}

// NewListHandler は ListHandler を生成する。
func NewListHandler(uc *user.ListUsersUsecase, logger *slog.Logger) *ListHandler {
	return &ListHandler{
		uc:     uc,
		logger: logger,
	}
}

// ServeHTTP はユーザー一覧を取得する。
// GET /users
func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	output, err := h.uc.Execute(r.Context())
	if err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(newListUsersResponse(output))
}
