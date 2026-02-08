package di

import (
	usecase "go-api/internal/application/user"
	"go-api/internal/infrastructure/repository/postgres"
	userhandler "go-api/internal/presentation/http/handler/user"
	sqlcuser "go-api/internal/sqlc/user"
)

// ListUserHandler はユーザー一覧取得ハンドラーを生成する。
func (c *Container) ListUserHandler() *userhandler.ListHandler {
	queries := sqlcuser.New(c.pool)
	repo := postgres.NewUserRepository(queries)
	uc := usecase.NewListUsersUsecase(repo)
	return userhandler.NewListHandler(uc, c.logger)
}

// CreateUserHandler はユーザー作成ハンドラーを生成する。
func (c *Container) CreateUserHandler() *userhandler.CreateHandler {
	queries := sqlcuser.New(c.pool)
	repo := postgres.NewUserRepository(queries)
	uc := usecase.NewCreateUserUsecase(repo)
	return userhandler.NewCreateHandler(uc, c.logger)
}
