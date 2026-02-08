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
	listUsersUC := usecase.NewListUsersUsecase(repo)
	createUserUC := usecase.NewCreateUserUsecase(repo)
	return handler.NewUserHandler(listUsersUC, createUserUC, c.logger)
}
