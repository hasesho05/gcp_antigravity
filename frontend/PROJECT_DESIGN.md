# Frontend Project Design

このドキュメントは、GCP認定資格模擬試験プラットフォームのフロントエンド開発における設計方針、技術スタック、アーキテクチャパターン、および開発ワークフローについて記述したものです。本設計は、Bulletproof React のアーキテクチャパターンをベースに、スケーラビリティと保守性の高いコードベースを目指しています。

## 目次

1.  [Tech Stack](#1-tech-stack)
2.  [Directory Structure (Bulletproof React)](#2-directory-structure-bulletproof-react)
3.  [Architecture & Patterns](#3-architecture--patterns)
    3.1. [API Type Integration (Quicktype)](#31-api-type-integration-quicktype)
    3.2. [Feature-based Architecture](#32-feature-based-architecture)
    3.3. [API Communication (useSWR)](#33-api-communication-useswr)
    3.4. [State Management](#34-state-management)
4.  [UI/Component Design](#4-uicomponent-design)
    4.1. [shadcn/ui & TailwindCSS](#41-shadcnui--tailwindcss)
    4.2. [Component Responsibility](#42-component-responsibility)
5.  [Testing Strategy](#5-testing-strategy)
    5.1. [Unit Testing (Vitest)](#51-unit-testing-vitest)
    5.2. [Component Testing (Vitest + Testing Library)](#52-component-testing-vitest--testing-library)
    5.3. [Storybook](#53-storybook)
6.  [Development Workflow](#6-development-workflow)
    6.1. [Code Quality](#61-code-quality)
    6.2. [Type Synchronization](#62-type-synchronization)
    6.3. [Adding a Feature](#63-adding-a-feature)

---

GCP認定資格模擬試験プラットフォームのフロントエンド設計書です。Bulletproof React のアーキテクチャパターンを採用し、スケーラビリティと保守性の高いコードベースを目指します。1. Tech StackCategoryTechnologyUsageFrameworkReact 18+UI LibraryBuild ToolVite高速な開発サーバーとビルドLanguageTypeScript型安全性 (Backend型との連携)StylingTailwindCSSUtility-first CSSUI Componentsshadcn/uiRadix UIベースの再利用可能なコンポーネントLinter/FormatterBiome高速なLint/Format (ESLint/Prettierの代替)TestingVitestユニットテスト・コンポーネントテストStorybookStorybookUIコンポーネントのカタログ化・開発Data FetchinguseSWRサーバー状態管理・データフェッチ・キャッシュGlobal StateZustandクライアント状態管理 (試験中の進行状態など)RoutingReact Router v6クライアントサイドルーティング2. Directory Structure (Bulletproof React)機能（Feature）ごとにディレクトリを分割する構造を採用します。src/
├── assets/              # 画像、フォントなどの静的アセット
├── components/          # アプリケーション全体で共有するUIコンポーネント
│   ├── ui/              # shadcn/ui によって生成された基本コンポーネント (Button, Input等)
│   ├── elements/        # shadcn以外の共有原子コンポーネント (LoadingSpinner等)
│   └── layouts/         # ページレイアウト (MainLayout, AuthLayout)
├── config/              # 環境変数やグローバル設定 (env.ts)
├── features/            # 機能単位のモジュール (ドメインロジックの核)
│   ├── auth/            # 認証機能 (Login, Register)
│   ├── exam/            # 試験機能 (問題表示、回答、結果)
│   ├── dashboard/       # ダッシュボード (成績表示、履歴)
│   └── admin/           # 管理者機能 (問題入稿)
│       ├── api/         # この機能固有のAPI取得関数 (Fetcher & SWR Hooks)
│       ├── components/  # この機能固有のコンポーネント
│       ├── hooks/       # この機能固有のカスタムフック (UIロジック)
│       ├── routes/      # この機能のルート定義
│       ├── stores/      # この機能固有の状態管理 (Zustand slices)
│       ├── types/       # この機能固有の型定義 (API型を拡張する場合など)
│       └── index.ts     # 公開APIのエントリーポイント
├── hooks/               # 汎用カスタムフック (useDisclosure等)
├── lib/                 # ライブラリの設定・ラッパー (axios, utils, swrConfig)
├── providers/           # アプリケーションプロバイダー (AppProvider)
├── routes/              # ルーティング設定 (AppRoutes)
├── stores/              # グローバルなクライアント状態 (UserStoreなど)
├── test/                # テストセットアップ、モック
├── types/               # **自動生成されたAPI型定義 (api.ts)**
└── utils/               # 純粋なユーティリティ関数
3. Architecture & Patterns3.1. API Type Integration (Quicktype)バックエンド（Go）の構造体から生成されたJSONを元に、TypeScriptの型定義を自動生成します。Source: backend/frontend_types_sample.json (Makefileで生成)Destination: frontend/src/types/api.tsRule: src/types/api.ts は手動で編集してはいけません。バックエンドの変更に合わせて再生成します。// src/features/exam/types/index.ts の例
import { Question } from '@/types/api';

// API型をUI用に拡張する場合のみ、Feature内で定義する
export type QuestionWithStatus = Question & {
  isAnswered: boolean;
};
3.2. Feature-based Architecture機能に関連するコードはすべて features/ 下に集約します。他のFeatureからインポートする場合は、必ず index.ts を経由し、内部構造への直接アクセス（Deep Import）を避けます。良い例: import { ExamList } from '@/features/exam';悪い例: import { ExamList } from '@/features/exam/components/ExamList';3.3. API Communication (useSWR)Client: axios を src/lib/axios.ts で設定し、InterceptorでFirebase AuthのTokenを付与します。Data Fetching: useSWR を使用します。useEffect 内でのデータフェッチは禁止します。features/{feature}/api/ 内に fetcher 関数と、useSWR をラップしたカスタムフックを作成します。// features/exam/api/getQuestions.ts
import useSWR from 'swr';
import { axios } from '@/lib/axios';
import { Question } from '@/types/api';

// Fetcher関数 (Axiosを利用)
const getQuestionsFetcher = (url: string): Promise<Question[]> => {
  return axios.get(url).then((res) => res.data);
};

// Custom Hook
export const useQuestions = (examId: string, setId: string) => {
  const { data, error, isLoading, mutate } = useSWR<Question[]>(
    `/exams/${examId}/sets/${setId}/questions`,
    getQuestionsFetcher
  );

  return {
    questions: data,
    isLoading,
    isError: error,
    mutate,
  };
};
3.4. State Management状態の種類に応じて適切な手法を選択します。Server State (API Cache): useSWRAPIからのデータは必ずSWR経由で取得し、キャッシュを利用します。Global Client State (User Session, Theme): ZustandFeature Local State (Exam Progress): Zustand試験中の回答状況などは、リロード対策として persist ミドルウェアを用いて localStorage と同期することを推奨します。Form State: React Hook Form + zod4. UI/Component Design4.1. shadcn/ui & TailwindCSS基本コンポーネントは src/components/ui に配置されます（npx shadcn-ui@latest add [component] で追加）。スタイリングは TailwindCSS のユーティリティクラスを使用します。複雑なスタイル分岐には clsx または tailwind-merge (cn ユーティリティ) を使用します。4.2. Component ResponsibilityContainer/Page Components: features/{feature}/routes/ に配置。useSWR カスタムフックを呼び出し、データを取得する責務を持ちます。ローディング状態 (isLoading) やエラー状態 (isError) のハンドリングを行います。Presentational Components: features/{feature}/components/ に配置。Propsを受け取って表示するだけ（Storybookの対象）。内部でAPIコールを行ってはいけません。5. Testing Strategy5.1. Unit Testing (Vitest)utils/ 内のロジック関数。hooks/ 内のカスタムフック。5.2. Component Testing (Vitest + Testing Library)複雑なインタラクションを持つUIコンポーネント。APIモック（MSW等）を使用したFeatureコンポーネントの統合テスト。5.3. Storybooksrc/components/ui の共通コンポーネント。features/{feature}/components のプレゼンテーションコンポーネント。各コンポーネントの stories.tsx を作成し、カタログとして管理します。6. Development Workflow6.1. Code QualityBiome を使用して Lint と Format を実行します。コミット前に npm run lint が通ることを確認します。6.2. Type Synchronizationバックエンドに変更があった場合：Backend: make generate-sampleFrontend: npm run gen:types (Quicktypeを実行するスクリプト)6.3. Adding a Featuresrc/features/ にディレクトリを作成。必要なAPI型が src/types/api.ts にあるか確認。api/ で useSWR ラッパーフックを作成。components/ でUIを作成 (Storybookで確認)。routes/ でページコンポーネントを作成し、Hooksを接続。src/routes/ にルートを追加。

---

この設計書は、GCP認定資格プラットフォームのフロントエンド開発において、一貫性のある高品質な開発を推進するための指針となります。ここに記載された原則とパターンに従うことで、効率的かつ保守性の高いアプリケーションを構築できると確信しています。