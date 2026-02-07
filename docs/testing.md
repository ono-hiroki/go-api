# Go テスト実行ガイド

## 基本コマンド

```bash
# 全テスト実行
go test ./...

# 特定パッケージのテスト
go test ./app/domain/user/valueobject

# 詳細出力
go test ./... -v
```

## パッケージ指定

| 指定方法 | 対象 |
|---------|------|
| `./...` | カレントディレクトリ以下すべて |
| `./app/...` | app ディレクトリ以下すべて |
| `./app/domain/user` | 特定パッケージのみ |
| なし | カレントディレクトリのみ |

## -run によるフィルタリング

`-run` は正規表現でテスト関数名をフィルタします。

```bash
# "UserID" を含むテスト
go test ./... -run "UserID" -v

# 前方一致
go test ./... -run "^TestNew" -v

# 完全一致
go test ./... -run "^TestNewUserID$" -v

# OR条件
go test ./... -run "UserID|Email" -v

# サブテストを指定
go test ./... -run "TestNewUserID/正常系" -v
```

## キャッシュ制御

```bash
# キャッシュを無視して実行
go test ./... -count=1

# キャッシュをクリア
go clean -testcache
```

`(cached)` と表示されたらキャッシュが使われています。

## よく使うオプション

| オプション | 意味 |
|-----------|------|
| `-v` | 詳細出力（テスト名を表示） |
| `-run "regex"` | 実行するテストをフィルタ |
| `-count=1` | キャッシュを無視 |
| `-short` | 短縮モード（`testing.Short()` で判定可能） |
| `-timeout 30s` | タイムアウト指定 |
| `-cover` | カバレッジ表示 |
| `-race` | データ競合検出 |

## gotestsum を使う場合

```bash
# 基本
gotestsum -- ./...

# フォーマット指定
gotestsum --format testname -- ./...

# フィルタ付き
gotestsum -- ./app/domain/user/valueobject -run "UserID" -v

# キャッシュ無視
gotestsum -- ./... -count=1
```

`gotestsum` は `--` の後に `go test` の引数を渡します。

## Taskfile を使う場合

```bash
# 全テスト
task test

# フィルタ付き
task test -- -run TestNewEmail/異常系

# 簡易出力
task test:short
```

## テスト内のログ出力

```go
func TestExample(t *testing.T) {
    // t.Log は -v 付きで表示される
    t.Log("デバッグ情報")
    t.Logf("値: %+v", obj)

    // fmt.Println は go test -v では表示されるが
    // gotestsum では抑制される場合がある
    fmt.Println("これは表示されないことがある")
}
```

テストでは `t.Log()` を使うのが推奨です。
