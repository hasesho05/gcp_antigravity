GCP Certification Exam Platform (Personal Project)

個人開発によるGCP認定資格の模擬試験プラットフォームです。
Udemyの模試形式を参考に、Webアプリとして低コストで運用し、将来的なネイティブアプリ化も見据えた設計になっています。

🚀 Tech Stack

Category

Technology

Note

Frontend

React (TypeScript, Vite)

SPA, Cloudflare Pages (想定)

Backend

Go (1.22+)

Cloud Run, Standard net/http

Database

Firestore

NoSQL, Cost-optimized schema

Auth

Firebase Authentication



Infra/CDN

Cloudflare

DNS, CDN, Frontend Hosting

Tools

Quicktype

Go構造体からTS型定義を自動生成

🏗 Architecture

Clean Architecture + Domain Driven Design (Lightweight)

Backend:

handler: HTTPリクエストの受付・レスポンス返却

usecase: アプリケーション固有のビジネスロジック。Input/Output DTOを定義。

domain: エンタープライズビジネスルール（エンティティ定義）

repository: データアクセスのインターフェース定義

repository_impl: リポジトリの具象実装。ドメインモデルとDBモデルの変換ロジックを持つ。

infra: Firestoreクライアントの初期化や共通処理などの低レイヤー実装。

Error Handling: github.com/cockroachdb/errors を使用し、スタックトレース付きのエラー管理を行う。

📂 Directory Structure

.
├── Makefile                 # ビルド、テスト、型定義生成コマンド
├── cmd
│   └── api
│       └── main.go          # エントリーポイント
├── internal
│   ├── domain               # ドメイン層 (純粋なエンティティ)
│   │   ├── question.go      # Question, AnswerOption
│   │   ├── attempt.go       # Attempt, AttemptStatus
│   │   ├── stats.go         # UserExamStats, DomainScore
│   │   └── error.go         # Domain Errors
│   ├── handler              # プレゼンテーション層
│   │   ├── admin            # 管理者用ハンドラ
│   │   │   └── handler.go
│   │   └── client           # クライアント用ハンドラ
│   │       └── handler.go
│   ├── usecase              # ユースケース層
│   │   ├── exam.go          # UseCase実装 (Interactor)
│   │   ├── input            # Input DTO (Request)
│   │   │   └── exam.go
│   │   └── output           # Output DTO (Response)
│   │       └── exam.go
│   ├── repository           # リポジトリインターフェース
│   │   └── repository.go
│   ├── repository_impl      # リポジトリ実装層
│   │   └── exam.go          # ExamRepositoryの実装 (Firestore操作ロジック)
│   └── infra                # インフラ層 (ドライバ/クライアント)
│       └── firestore        # Firestore共通処理
│           └── client.go    # Client初期化、共通ヘルパー
└── scripts
    └── dump_json.go         # Quicktype用JSON生成スクリプト


🛠 Development Commands

1. Run Backend

make run


2. Generate TypeScript Types (for Frontend)

GoのInput/Output DTOおよびDomain定義からJSONサンプルを出力し、それを元にフロントエンド用の型定義を作成します。

# 1. JSONサンプルを出力
make generate-sample > frontend_types_sample.json

# 2. (Optional) Quicktype CLIでTS型を生成
quicktype -o frontend/src/types/api.ts --src frontend_types_sample.json --src-lang json --lang ts


3. Test

make test
