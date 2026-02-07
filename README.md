# go-api

Go言語によるREST APIサーバー。クリーンアーキテクチャで構築。

## 必要なツール

- Go 1.25+
- Docker / Docker Compose
- [mise](https://mise.jdx.dev/)（推奨）
- [Task](https://taskfile.dev/)

## セットアップ

```bash
# ツールのインストール
mise install

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

# フォーマット + 静的解析 + テスト
task check
```

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
├── app/
│   └── domain/              # ドメイン層
│       └── user/
│           ├── user.go      # Userエンティティ
│           ├── repository.go
│           └── valueobject/
├── db/
│   └── migrations/          # マイグレーションファイル
├── docs/                    # 開発ドキュメント
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
| マイグレーション | golang-migrate |
| アーキテクチャ | クリーンアーキテクチャ |
