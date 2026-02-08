package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go-api/internal/application/user"
	"go-api/internal/domain"
	httperrors "go-api/internal/presentation/http/errors"
	"go-api/internal/presentation/http/validation"
)

// updateUserRequest はユーザー更新のJSONリクエスト。
type updateUserRequest struct {
	Name  string `json:"name" validate:"required,max=100"`
	Email string `json:"email" validate:"required,max=255,email"`
}

func (r updateUserRequest) toInput() user.UpdateUserInput {
	return user.UpdateUserInput{
		Name:  r.Name,
		Email: r.Email,
	}
}

// updateUserResponse はユーザー更新のJSONレスポンス。
type updateUserResponse struct {
	User updateUserResponseUser `json:"user"`
}

type updateUserResponseUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func newUpdateUserResponse(output *user.UpdateUserOutput) updateUserResponse {
	return updateUserResponse{
		User: updateUserResponseUser{
			ID:    output.User.ID,
			Name:  output.User.Name,
			Email: output.User.Email,
		},
	}
}

// UpdateHandler はユーザー更新のHTTPハンドラー。
type UpdateHandler struct {
	uc     *user.UpdateUserUsecase
	logger *slog.Logger
}

// NewUpdateHandler は UpdateHandler を生成する。
func NewUpdateHandler(uc *user.UpdateUserUsecase, logger *slog.Logger) *UpdateHandler {
	return &UpdateHandler{
		uc:     uc,
		logger: logger,
	}
}

// ServeHTTP はユーザーを更新する。
// PUT /users/{id}
func (h *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperrors.WriteError(w, r, domain.ErrInvalidInput, h.logger)
		return
	}

	if err := validation.Struct(req); err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	output, err := h.uc.Execute(r.Context(), id, req.toInput())
	if err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(newUpdateUserResponse(output))
}
