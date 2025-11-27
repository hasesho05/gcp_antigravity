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

Tygo

Go構造体からTS型定義を自動生成 (Union型対応)

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
│   ├── api
│   │   └── main.go          # エントリーポイント
│   └── seed_questions       # 問題データ投入スクリプト
│       ├── main.go
│       └── source.json
├── internal
│   ├── domain               # ドメイン層 (純粋なエンティティ)
│   │   ├── question.go      # Question, AnswerOption
│   │   ├── exam_set.go      # ExamSet
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
├── tygo.yaml                # Tygo設定ファイル

## 🧩 Domain Models

アプリケーションの中核となるドメインモデルの解説です。

### 1. Exam (資格試験)
GCP認定試験そのものを表します（例: "Professional Cloud Developer"）。
- `ID`: 資格ID (e.g. `professional-cloud-developer`)
- `Code`: 資格コード (e.g. `PCD`)

### 2. ExamSet (模擬試験セット)
1つの資格試験に含まれる、50問1セットの模擬試験単位です。
- `ID`: セットID (e.g. `practice_exam_1`)
- `ExamID`: 親となる資格ID

### 3. Question (問題)
個々の問題データです。
- `ID`: 問題ID (e.g. `PCD_SET1_001`)
- `QuestionType`: `multiple-choice` (単一選択) または `multi-select` (複数選択)
- `Domain`: 出題分野 (e.g. "Identity and Security")
- `OverallExplanation`: 全体の解説

### 4. User (ユーザー)
アプリケーションの利用者です。Firebase Authと連携します。
- `Role`: `free` (無料), `pro` (有料), `admin` (管理者)
- `SubscriptionStatus`: サブスクリプションの状態

### 5. Attempt (受験)
ユーザーが模擬試験を1回受験した履歴を表します。
- `Status`: `in_progress` (受験中), `paused` (中断), `completed` (完了)
- `Answers`: ユーザーの回答状況 (Map形式)
- `Score`: 正解数

### 6. UserExamStats (成績統計)
ユーザーの資格ごとの累積成績です。
- `DomainStats`: 分野ごとの正答率などの統計情報
- `AccuracyRate`: 全体の正答率


🛠 Development Commands

1. Run Backend

make run


2. Generate TypeScript Types (for Frontend)

Goのドメイン定義からTypeScriptの型定義を自動生成します（Tygo使用）。

```bash
make generate-types
```

`frontend/src/types/api.ts` が更新されます。


3. Test
make test

4. Seed Questions (Development)

JSONファイルから問題データをFirestoreに投入します。

```bash
# 1. backend/cmd/seed_questions/source.json に問題データを配置
# 2. 以下のコマンドを実行
make seed-questions
```

## 🔥 Firestore Data Structure

```
exams/{examID}
  ├── sets/{setID} (ExamSet)
  │     ├── questions/{questionID} (Question)
```

- **ExamSet**: 模擬試験のセット（例: "Practice Exam 1"）
- **Question**: 個々の問題データ
