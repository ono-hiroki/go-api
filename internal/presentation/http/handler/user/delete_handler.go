package user

import (
	"log/slog"
	"net/http"

	"go-api/internal/application/user"
	httperrors "go-api/internal/presentation/http/errors"
)

// DeleteHandler はユーザー削除のHTTPハンドラー。
type DeleteHandler struct {
	uc     *user.DeleteUserUsecase
	logger *slog.Logger
}

// NewDeleteHandler は DeleteHandler を生成する。
func NewDeleteHandler(uc *user.DeleteUserUsecase, logger *slog.Logger) *DeleteHandler {
	return &DeleteHandler{
		uc:     uc,
		logger: logger,
	}
}

// ServeHTTP はユーザーを削除する。
// DELETE /users/{id}
func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.uc.Execute(r.Context(), id); err != nil {
		httperrors.WriteError(w, r, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
