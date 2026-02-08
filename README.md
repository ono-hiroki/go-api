# go-api

Go言語によるREST APIサーバー。クリーンアーキテクチャで構築。

## 必要なツール

- Go 1.25+
- Node.js 22+ (TypeSpec用)
- Docker / Docker Compose
- [mise](https://mise.jdx.dev/)（推奨）
- [Task](https://taskfile.dev/)

## セットアップ

```bash
# ツールのインストール
mise install

# npm依存関係インストール
npm install

# PostgreSQL起動
docker compose up -d

# マイグレーション実行
task db:up
```

## 開発コマンド

```bash
# サーバー起動
task run

# テスト実行
task test

# テスト（簡易出力）
task test:short

# インテグレーションテスト
task test:integration

# フォーマット + 静的解析 + テスト
task check

# OpenAPI生成 (TypeSpecから)
npm run openapi
```

## API エンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| GET | /health | ヘルスチェック |
| GET | /users | ユーザー一覧取得 |
| POST | /users | ユーザー作成 |
| GET | /users/{id} | ユーザー取得 |
| PUT | /users/{id} | ユーザー更新 |
| DELETE | /users/{id} | ユーザー削除 |

API仕様の詳細は [api/openapi.yaml](api/openapi.yaml) を参照。

## DB操作

```bash
# マイグレーション実行
task db:up

# ロールバック（1つ戻す）
task db:down

# DBリセット（全削除→再作成）
task db:reset

# 新規マイグレーション作成
task db:create -- create_posts
```

## ディレクトリ構成

```
go-api/
├── cmd/api/                 # エントリーポイント
├── internal/
│   ├── application/         # ユースケース層
│   │   └── user/
│   ├── domain/              # ドメイン層
│   │   └── user/
│   │       └── valueobject/
│   ├── infrastructure/      # インフラ層
│   │   └── repository/
│   │       └── postgres/
│   ├── presentation/        # プレゼンテーション層
│   │   └── http/
│   │       ├── handler/user/
│   │       ├── errors/
│   │       ├── middleware/
│   │       └── validation/
│   └── di/                  # 依存性注入
├── api/                     # 生成されたOpenAPI
├── typespec/                # TypeSpec定義
├── db/migrations/           # マイグレーションファイル
├── docker-compose.yml
├── Taskfile.yml
└── mise.toml
```

## 技術スタック

| カテゴリ | 技術 |
|---------|------|
| 言語 | Go 1.25 |
| HTTP | 標準ライブラリ (net/http) |
| DB | PostgreSQL 17 |
| ORM/クエリ | sqlc |
| マイグレーション | golang-migrate |
| バリデーション | go-playground/validator |
| API定義 | TypeSpec → OpenAPI 3.0 |
| アーキテクチャ | クリーンアーキテクチャ |
