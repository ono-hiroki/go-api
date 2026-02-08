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

// createUserRequest はユーザー作成のJSONリクエスト。
type createUserRequest struct {
	Name  string `json:"name" validate:"required,max=100"`
	Email string `json:"email" validate:"required,max=255,email"`
}

func (r createUserRequest) toInput() user.CreateUserInput {
	return user.CreateUserInput{
		Name:  r.Name,
		Email: r.Email,
	}
}

// createUserResponse はユーザー作成のJSONレスポンス。
type createUserResponse struct {
	User createUserResponseUser `json:"user"`
}

type createUserResponseUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func newCreateUserResponse(output *user.CreateUserOutput) createUserResponse {
	return createUserResponse{
		User: createUserResponseUser{
			ID:    output.User.ID,
			Name:  output.User.Name,
			Email: output.User.Email,
		},
	}
}

// CreateHandler はユーザー作成のHTTPハンドラー。
type CreateHandler struct {
	uc     *user.CreateUserUsecase
	logger *slog.Logger
}

// NewCreateHandler は CreateHandler を生成する。
func NewCreateHandler(uc *user.CreateUserUsecase, logger *slog.Logger) *CreateHandler {
	return &CreateHandler{
		uc:     uc,
		logger: logger,
	}
}

// ServeHTTP はユーザーを作成する。
// POST /users
func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperrors.WriteError(w, r, domain.ErrInvalidInput, h.logger)
		return
	}

	if err := validation.Struct(req); err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	output, err := h.uc.Execute(r.Context(), req.toInput())
	if err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newCreateUserResponse(output))
}
