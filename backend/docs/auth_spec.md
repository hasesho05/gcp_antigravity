# ユーザー認証仕様書

## 概要
本システムでは、認証基盤として **Firebase Authentication** を採用しています。
クライアントサイドでFirebase Authを使用してIDトークンを取得し、バックエンドAPIへのリクエスト時にAuthorizationヘッダーにBearerトークンとして付与することで認証を行います。

## 認証フロー

1.  **クライアントサイド**:
    *   Firebase SDKを使用してユーザー登録・ログインを行う（Googleログイン、メール/パスワード等）。
    *   ログイン成功後、IDトークンを取得する。
    *   APIリクエストのヘッダーに `Authorization: Bearer <ID_TOKEN>` を付与する。

2.  **バックエンドサイド**:
    *   **Auth Middleware**:
        *   リクエストヘッダーからIDトークンを検証する。
        *   検証に成功した場合、トークンからUID（User ID）を抽出し、Contextにセットする。
    *   **User Handler**:
        *   ContextからUIDを取得し、リクエストされた操作を行う。

## ユーザー管理

Firebase Authのユーザーとは別に、アプリケーション独自のユーザー情報をFirestoreの `users` コレクションで管理します。

### データモデル

| フィールド名 | 型 | 説明 |
| :--- | :--- | :--- |
| `id` | string | Firebase AuthのUID（ドキュメントIDとしても使用） |
| `email` | string | メールアドレス |
| `provider` | enum | 認証プロバイダー (`google.com`, `password`, `github.com`) |
| `role` | enum | ユーザー権限 (`free`, `pro`, `admin`) |
| `subscriptionStatus` | enum | サブスクリプション状態 (`active`, `expired`, `canceled`) |
| `createdAt` | timestamp | 作成日時 |

### API エンドポイント

#### 1. ユーザー作成 (`POST /users`)

Firebase Authでの登録完了後、バックエンドにユーザー情報を作成するために呼び出します。

*   **認証**: 必須 (Firebase ID Token)
*   **リクエストボディ**:
    ```json
    {
      "email": "user@example.com",
      "provider": "google.com" // オプション (デフォルト: "password")
    }
    ```
*   **レスポンス**:
    *   `201 Created`: 作成成功
    *   `409 Conflict`: メールアドレスが既に登録されている場合
    *   `400 Bad Request`: リクエスト不正

#### 2. 現在のユーザー取得 (`GET /users/me`)

ログイン中のユーザー自身の情報を取得します。

*   **認証**: 必須 (Firebase ID Token)
*   **レスポンス**:
    *   `200 OK`: ユーザー情報
    *   `404 Not Found`: ユーザー情報が存在しない場合

## バリデーション

*   **メールアドレスの重複**:
    *   `POST /users` 実行時に、指定されたメールアドレスが既に `users` コレクションに存在する場合、`409 Conflict` エラーを返します。
    *   これは、異なる認証プロバイダーであってもメールアドレスの一意性を保つための仕様です。

## 将来の拡張性

`provider` フィールドにより、Google認証以外のプロバイダー（GitHub, Twitter等）や、メール/パスワード認証など、複数の認証方式を識別・管理できるように設計されています。
