package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-api/internal/application/user"
	httperrors "go-api/internal/presentation/http/errors"
)

// getUserResponse はユーザー取得のJSONレスポンス。
type getUserResponse struct {
	User getUserResponseUser `json:"user"`
}

type getUserResponseUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func newGetUserResponse(output *user.GetUserOutput) getUserResponse {
	return getUserResponse{
		User: getUserResponseUser{
			ID:    output.User.ID,
			Name:  output.User.Name,
			Email: output.User.Email,
		},
	}
}

// GetHandler はユーザー取得のHTTPハンドラー。
type GetHandler struct {
	uc     *user.GetUserUsecase
	logger *slog.Logger
}

// NewGetHandler は GetHandler を生成する。
func NewGetHandler(uc *user.GetUserUsecase, logger *slog.Logger) *GetHandler {
	return &GetHandler{
		uc:     uc,
		logger: logger,
	}
}

// ServeHTTP はユーザーを取得する。
// GET /users/{id}
func (h *GetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	output, err := h.uc.Execute(r.Context(), id)
	if err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(newGetUserResponse(output))
}
