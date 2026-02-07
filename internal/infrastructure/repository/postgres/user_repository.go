// Package postgres はPostgreSQLを使用したリポジトリ実装を提供する。
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"go-api/internal/domain"
	"go-api/internal/domain/user"
	"go-api/internal/domain/user/valueobject"
	sqlcuser "go-api/internal/sqlc/user"
)

const (
	// PostgreSQL エラーコード
	pgUniqueViolation = "23505"
)

// UserRepository はPostgreSQLを使用したユーザーリポジトリの実装。
type UserRepository struct {
	queries *sqlcuser.Queries
}

// NewUserRepository は UserRepository を生成する。
func NewUserRepository(queries *sqlcuser.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

// Save はユーザーをDBに保存する。
// 一意制約違反の場合は domain.ErrConflict を返す。
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	err := r.queries.CreateUser(ctx, sqlcuser.CreateUserParams{
		ID:    uuidToPgtype(u.ID()),
		Name:  u.Name().String(),
		Email: u.Email().String(),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.Conflict("user", "Save", err)
		}
		return err
	}
	return nil
}

// FindByID は指定されたIDのユーザーを取得する。
// 見つからない場合は domain.ErrNotFound を返す。
func (r *UserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*user.User, error) {
	row, err := r.queries.GetUser(ctx, uuidToPgtype(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NotFound("user", "FindByID")
		}
		return nil, err
	}
	return toEntity(&row)
}

// FindAll は全ユーザーを取得する。
func (r *UserRepository) FindAll(ctx context.Context) ([]*user.User, error) {
	rows, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*user.User, 0, len(rows))
	for i := range rows {
		u, err := toEntity(&rows[i])
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Delete は指定されたIDのユーザーを削除する。
func (r *UserRepository) Delete(ctx context.Context, id valueobject.UserID) error {
	return r.queries.DeleteUser(ctx, uuidToPgtype(id))
}

// uuidToPgtype はドメインのUserIDをPostgreSQLのUUID型に変換する。
func uuidToPgtype(id valueobject.UserID) pgtype.UUID {
	var pgID pgtype.UUID
	_ = pgID.Scan(id.String())
	return pgID
}

// toEntity はsqlcの行データをドメインのUserエンティティに変換する。
func toEntity(row *sqlcuser.User) (*user.User, error) {
	id, err := valueobject.ParseUserID(uuidToString(row.ID))
	if err != nil {
		return nil, err
	}
	name, err := valueobject.NewUserName(row.Name)
	if err != nil {
		return nil, err
	}
	email, err := valueobject.NewEmail(row.Email)
	if err != nil {
		return nil, err
	}
	return user.Reconstruct(id, name, email), nil
}

// uuidToString はPostgreSQLのUUID型を文字列に変換する。
func uuidToString(id pgtype.UUID) string {
	if !id.Valid {
		return ""
	}
	b := id.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
