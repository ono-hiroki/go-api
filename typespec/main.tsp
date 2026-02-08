import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;

@service(#{
  title: "Go API",
})
@server("http://localhost:8080", "Local development server")
namespace GoApi;

// ========================================
// Common Models
// ========================================

/** バリデーションエラーの詳細 */
model ValidationErrorDetail {
  /** エラーが発生したフィールド名 */
  field: string;

  /** エラーコード (required, too_long, invalid_format など) */
  code: string;

  /** エラーメッセージ */
  message: string;
}

/** バリデーションエラーレスポンス */
@error
model ValidationError {
  @statusCode statusCode: 400;
  @body body: {
    code: "VALIDATION_ERROR";
    message: string;
    details: ValidationErrorDetail[];
  };
}

/** 内部サーバーエラーレスポンス */
@error
model InternalServerError {
  @statusCode statusCode: 500;
  @body body: {
    code: "INTERNAL_ERROR";
    message: "internal server error";
  };
}

/** リソースが見つからないエラー */
@error
model NotFoundError {
  @statusCode statusCode: 404;
  @body body: {
    code: "NOT_FOUND";
    message: string;
  };
}

// ========================================
// User Models
// ========================================

/** ユーザー情報 */
model User {
  /** ユーザーID (UUID) */
  id: string;

  /** ユーザー名 */
  name: string;

  /** メールアドレス */
  email: string;
}

/** ユーザー作成リクエスト */
model CreateUserRequest {
  /** ユーザー名 (1-100文字) */
  @maxLength(100)
  name: string;

  /** メールアドレス */
  @maxLength(255)
  email: string;
}

/** ユーザー一覧レスポンス */
model ListUsersResponse {
  users: User[];
}

/** ユーザー作成レスポンス */
model CreateUserResponse {
  user: User;
}

/** ユーザー取得レスポンス */
model GetUserResponse {
  user: User;
}

/** ユーザー更新リクエスト */
model UpdateUserRequest {
  /** ユーザー名 (1-100文字) */
  @maxLength(100)
  name: string;

  /** メールアドレス */
  @maxLength(255)
  email: string;
}

/** ユーザー更新レスポンス */
model UpdateUserResponse {
  user: User;
}

// ========================================
// User API
// ========================================

@route("/users")
@tag("Users")
interface Users {
  /** ユーザー一覧を取得する */
  @get
  list(): ListUsersResponse | InternalServerError;

  /** ユーザーを作成する */
  @post
  create(@body body: CreateUserRequest): {
    @statusCode statusCode: 201;
    @body body: CreateUserResponse;
  } | ValidationError | InternalServerError;

  /** ユーザーを取得する */
  @get
  @route("{id}")
  get(@path id: string): GetUserResponse | NotFoundError | InternalServerError;

  /** ユーザーを更新する */
  @put
  @route("{id}")
  update(@path id: string, @body body: UpdateUserRequest): UpdateUserResponse | ValidationError | NotFoundError | InternalServerError;

  /** ユーザーを削除する */
  @delete
  @route("{id}")
  delete(@path id: string): {
    @statusCode statusCode: 204;
  } | NotFoundError | InternalServerError;
}

// ========================================
// Health Check
// ========================================

@route("/health")
@tag("System")
interface Health {
  /** ヘルスチェック */
  @get
  check(): {
    status: "ok";
  };
}
