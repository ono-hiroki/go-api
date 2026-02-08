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

// GetUserHandler はユーザー取得ハンドラーを生成する。
func (c *Container) GetUserHandler() *userhandler.GetHandler {
	queries := sqlcuser.New(c.pool)
	repo := postgres.NewUserRepository(queries)
	uc := usecase.NewGetUserUsecase(repo)
	return userhandler.NewGetHandler(uc, c.logger)
}

// UpdateUserHandler はユーザー更新ハンドラーを生成する。
func (c *Container) UpdateUserHandler() *userhandler.UpdateHandler {
	queries := sqlcuser.New(c.pool)
	repo := postgres.NewUserRepository(queries)
	uc := usecase.NewUpdateUserUsecase(repo)
	return userhandler.NewUpdateHandler(uc, c.logger)
}

// DeleteUserHandler はユーザー削除ハンドラーを生成する。
func (c *Container) DeleteUserHandler() *userhandler.DeleteHandler {
	queries := sqlcuser.New(c.pool)
	repo := postgres.NewUserRepository(queries)
	uc := usecase.NewDeleteUserUsecase(repo)
	return userhandler.NewDeleteHandler(uc, c.logger)
}
