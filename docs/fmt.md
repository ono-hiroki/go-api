# fmt パッケージ

## 関数一覧

### 出力先による分類

| 関数 | 出力先 |
|------|--------|
| `Print`, `Println`, `Printf` | 標準出力 |
| `Fprint`, `Fprintln`, `Fprintf` | 任意の `io.Writer` |
| `Sprint`, `Sprintln`, `Sprintf` | 文字列を返す |

### サフィックスの違い

| サフィックス | 意味 |
|-------------|------|
| なし (`Print`) | そのまま出力 |
| `ln` (`Println`) | 末尾に改行を追加 |
| `f` (`Printf`) | フォーマット指定子を使用 |

### 使用例

```go
name := "Go"

// Print系（標準出力）
fmt.Print(name)               // "Go"（改行なし）
fmt.Println(name)             // "Go\n"（改行あり）
fmt.Printf("Hello %s", name)  // "Hello Go"（フォーマット）

// Sprint系（文字列を返す）
s := fmt.Sprintf("Hello %s", name)  // s = "Hello Go"

// Fprint系（io.Writerに出力）
fmt.Fprintln(os.Stderr, "error!")   // 標準エラーに出力
fmt.Fprintf(w, "Hello %s", name)    // http.ResponseWriter等に
```

### よく使う場面

| 関数 | 用途 |
|------|------|
| `Println` | デバッグ出力 |
| `Printf` | フォーマット付きデバッグ |
| `Sprintf` | 文字列生成 |
| `Fprintf` | ファイル・HTTPレスポンス等への出力 |
| `Errorf` | エラー生成（`errors.New` の代わり） |

---

## フォーマット指定子（verb）

### 汎用

| verb | 意味 | 例 |
|------|------|-----|
| `%v` | デフォルト形式 | `{田中太郎 tanaka@example.com}` |
| `%+v` | フィールド名付き | `{name:田中太郎 email:tanaka@example.com}` |
| `%#v` | Go構文形式 | `user.User{name:"田中太郎", email:"..."}` |
| `%T` | 型名 | `user.User` |
| `%%` | リテラル% | `%` |

## 真偽値

| verb | 意味 | 例 |
|------|------|-----|
| `%t` | true または false | `true` |

## 整数

| verb | 意味 | 例 |
|------|------|-----|
| `%d` | 10進数 | `42` |
| `%b` | 2進数 | `101010` |
| `%o` | 8進数 | `52` |
| `%x` | 16進数（小文字） | `2a` |
| `%X` | 16進数（大文字） | `2A` |
| `%c` | Unicode文字 | `*`（42の場合） |
| `%U` | Unicodeコードポイント | `U+002A` |

## 浮動小数点

| verb | 意味 | 例 |
|------|------|-----|
| `%f` | 小数点表記 | `3.140000` |
| `%e` | 指数表記（小文字） | `3.140000e+00` |
| `%E` | 指数表記（大文字） | `3.140000E+00` |
| `%g` | 簡潔な表記（%e か %f） | `3.14` |

## 文字列・バイト列

| verb | 意味 | 例 |
|------|------|-----|
| `%s` | 文字列 | `hello` |
| `%q` | クォート付き文字列 | `"hello"` |
| `%x` | 16進数（各バイト） | `68656c6c6f` |

## ポインタ

| verb | 意味 | 例 |
|------|------|-----|
| `%p` | 16進数アドレス | `0xc0000b4000` |

## 幅・精度指定

```go
fmt.Printf("%10s", "go")      // "        go"  （幅10、右寄せ）
fmt.Printf("%-10s", "go")     // "go        "  （幅10、左寄せ）
fmt.Printf("%.2f", 3.14159)   // "3.14"        （小数点以下2桁）
fmt.Printf("%05d", 42)        // "00042"       （ゼロ埋め5桁）
fmt.Printf("%*d", 5, 42)      // "   42"       （幅を引数で指定）
```

## フラグ

| フラグ | 意味 |
|-------|------|
| `-` | 左寄せ |
| `+` | 符号を常に表示（+/-） |
| ` ` | 正の数の前にスペース |
| `0` | ゼロ埋め |
| `#` | 代替形式（0x, 0o 等を付加） |

```go
fmt.Printf("%+d", 42)   // "+42"
fmt.Printf("% d", 42)   // " 42"
fmt.Printf("%#x", 42)   // "0x2a"
fmt.Printf("%#o", 42)   // "0o52"
```

## デバッグでよく使うパターン

```go
// 構造体の中身確認
fmt.Printf("%+v\n", user)

// コピペで使える形式
fmt.Printf("%#v\n", user)

// 型と値を同時に表示
fmt.Printf("%T: %v\n", x, x)

// スライスの中身確認
fmt.Printf("items: %+v\n", items)
```

## 参考

- [公式ドキュメント](https://pkg.go.dev/fmt)
