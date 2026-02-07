//go:build integration

package postgres_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-api/internal/domain/user"
	"go-api/internal/domain/user/valueobject"
	"go-api/internal/infrastructure/repository/postgres"
	sqlcuser "go-api/internal/sqlc/user"
	"go-api/internal/testutil/factory"
)

const testTimeout = 5 * time.Second

var testPool *pgxpool.Pool

// TODO: testcontainers-go でフォールバック対応
// Docker未起動時に自動でコンテナ起動するようにする
func TestMain(m *testing.M) {
	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		log.Fatal("TEST_DATABASE_URL が設定されていません")
	}

	var err error
	testPool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("データベース接続に失敗: %v", err)
	}

	code := m.Run()

	testPool.Close()
	os.Exit(code)
}

func setupTest(t *testing.T) (context.Context, pgx.Tx, *postgres.UserRepository) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	t.Cleanup(cancel)

	tx, err := testPool.Begin(ctx)
	require.NoError(t, err, "トランザクション開始に失敗")

	t.Cleanup(func() {
		tx.Rollback(context.Background())
	})

	queries := sqlcuser.New(tx)
	return ctx, tx, postgres.NewUserRepository(queries)
}

// insertUserRow はDB直接INSERTでテストデータを作成する
func insertUserRow(t *testing.T, ctx context.Context, tx pgx.Tx, u *user.User) {
	t.Helper()
	_, err := tx.Exec(ctx,
		`INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`,
		u.ID().String(), u.Name().String(), u.Email().String(),
	)
	require.NoError(t, err, "テストデータのINSERTに失敗")
}

// selectUserRow はDB直接SELECTでユーザーを取得する
func selectUserRow(t *testing.T, ctx context.Context, tx pgx.Tx, id valueobject.UserID) (name, email string, found bool) {
	t.Helper()
	err := tx.QueryRow(ctx,
		`SELECT name, email FROM users WHERE id = $1`,
		id.String(),
	).Scan(&name, &email)
	if err == pgx.ErrNoRows {
		return "", "", false
	}
	require.NoError(t, err, "テストデータのSELECTに失敗")
	return name, email, true
}

func TestUserRepository_Save(t *testing.T) {
	t.Run("ユーザーを保存できる", func(t *testing.T) {
		ctx, tx, repo := setupTest(t)

		u := factory.NewUser(
			factory.WithName("保存テスト"),
			factory.WithEmail("save@example.com"),
		)

		err := repo.Save(ctx, u)
		require.NoError(t, err, "Save に失敗")

		// DB直読みで検証（repo.FindByIDを使わない）
		name, email, found := selectUserRow(t, ctx, tx, u.ID())
		require.True(t, found, "保存したユーザーがDBに存在しない")
		assert.Equal(t, u.Name().String(), name, "Name が一致しない")
		assert.Equal(t, u.Email().String(), email, "Email が一致しない")
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	t.Run("存在するIDでユーザーを取得できる", func(t *testing.T) {
		ctx, tx, repo := setupTest(t)

		u := factory.NewUser()
		insertUserRow(t, ctx, tx, u)

		found, err := repo.FindByID(ctx, u.ID())
		require.NoError(t, err, "FindByID に失敗")
		require.NotNil(t, found, "ユーザーが見つかるはずが nil")

		assert.Equal(t, u.ID().String(), found.ID().String(), "ID が一致しない")
		assert.Equal(t, u.Name().String(), found.Name().String(), "Name が一致しない")
		assert.Equal(t, u.Email().String(), found.Email().String(), "Email が一致しない")
	})

	t.Run("存在しないIDでnilを返す", func(t *testing.T) {
		ctx, _, repo := setupTest(t)

		nonExistentID := valueobject.NewUserID()

		found, err := repo.FindByID(ctx, nonExistentID)
		require.NoError(t, err, "FindByID に失敗")
		assert.Nil(t, found, "存在しないIDで nil が返るべき")
	})
}

func TestUserRepository_FindAll(t *testing.T) {
	t.Run("全ユーザーを取得できる", func(t *testing.T) {
		ctx, tx, repo := setupTest(t)

		u1 := factory.NewUser()
		u2 := factory.NewUser()
		insertUserRow(t, ctx, tx, u1)
		insertUserRow(t, ctx, tx, u2)

		users, err := repo.FindAll(ctx)
		require.NoError(t, err, "FindAll に失敗")
		assert.Len(t, users, 2, "2件取得できるべき")

		ids := []string{users[0].ID().String(), users[1].ID().String()}
		assert.Contains(t, ids, u1.ID().String(), "u1 が含まれるべき")
		assert.Contains(t, ids, u2.ID().String(), "u2 が含まれるべき")
	})
}

func TestUserRepository_Delete(t *testing.T) {
	t.Run("ユーザーを削除できる", func(t *testing.T) {
		ctx, tx, repo := setupTest(t)

		u := factory.NewUser()
		insertUserRow(t, ctx, tx, u)

		err := repo.Delete(ctx, u.ID())
		require.NoError(t, err, "Delete に失敗")

		// DB直読みで削除確認（repo.FindByIDを使わない）
		_, _, found := selectUserRow(t, ctx, tx, u.ID())
		assert.False(t, found, "削除後はDBに存在しないべき")
	})
}
