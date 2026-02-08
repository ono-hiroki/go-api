package di

import (
	usecase "go-api/internal/application/user"
	"go-api/internal/infrastructure/repository/postgres"
	"go-api/internal/presentation/http/handler"
	sqlcuser "go-api/internal/sqlc/user"
)

// UserHandler はUserHandlerを生成する。
func (c *Container) UserHandler() *handler.UserHandler {
	queries := sqlcuser.New(c.pool)
	repo := postgres.NewUserRepository(queries)
	uc := usecase.NewListUsersUsecase(repo)
	return handler.NewUserHandler(uc, c.logger)
}
