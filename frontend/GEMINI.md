# GCP認定資格プラットフォーム フロントエンド開発ガイドライン (GEMINI.md)

このドキュメントは、GCP認定資格模擬試験プラットフォームのReactフロントエンド開発における、アーキテクチャの原則、コーディング規則、および標準的な開発プラクティスを定義するものです。開発エージェント(Gemini)および人間の開発者は、コードベースの一貫性を保つため、このガイドラインに従う必要があります。

## 1. プロジェクト概要

このプロジェクトは、React + TypeScript + Viteを使用したGCP認定資格の模擬試験プラットフォームのフロントエンドアプリケーションです。

### 技術スタック

- **フレームワーク**: React 18 + TypeScript
- **ビルドツール**: Vite
- **スタイリング**: TailwindCSS
- **UIコンポーネント**: shadcn/ui
- **状態管理**: React Hooks
- **データフェッチング**: SWR
- **ルーティング**: TanStack Router
- **コード品質**: Biome (Linter & Formatter)
- **テスト**: Vitest + React Testing Library
- **コンポーネント開発**: Storybook

### プロジェクト構成

このプロジェクトは **Bulletproof React** のディレクトリ構造に基づいています。

```
src/
├── app/             # アプリケーションレイヤー（ルーティング、プロバイダー）
├── components/      # 共通UIコンポーネント（shadcn/ui）
├── features/        # 機能別モジュール（各機能のコンポーネント、API、型定義）
├── hooks/           # カスタムフック
├── lib/             # ライブラリ設定・ユーティリティ
├── stores/          # グローバル状態管理
├── types/           # グローバル型定義
└── utils/           # ユーティリティ関数
```

## 2. コーディング規約

### 2.1. TypeScript規約

#### 2.1.1 引数の受け取り方

**【禁止】Propsの分割代入（destructuring）**

```tsx
// ❌ 禁止: 引数の分割代入
const MyComponent = ({ title, description, isVisible }: {
  title: string;
  description: string;
  isVisible: boolean;
}) => {
  return (
    <div>
      <h1>{title}</h1>
      <p>{description}</p>
    </div>
  );
};

// ❌ 禁止: 型を別に定義しても分割代入は禁止
type MyComponentProps = {
  title: string;
  description: string;
  isVisible: boolean;
};

const MyComponent = ({ title, description, isVisible }: MyComponentProps) => {
  return <div>{title}</div>;
};
```

```tsx
// ✅ 推奨: propsオブジェクトとして受け取り、インライン型定義を使用
const MyComponent = (props: {
  title: string;
  description: string;
  isVisible: boolean;
}) => {
  return (
    <div>
      <h1>{props.title}</h1>
      <p>{props.description}</p>
    </div>
  );
};

// ✅ 推奨: 複雑な型は同じファイルの上部で定義
// components/ExamCard.tsx
type ExamCardProps = {
  exam: {
    id: string;
    name: string;
    description: string;
  };
  onClick: (id: string) => void;
};

export const ExamCard = (props: ExamCardProps) => {
  return (
    <div onClick={() => props.onClick(props.exam.id)}>
      <h3>{props.exam.name}</h3>
      <p>{props.exam.description}</p>
    </div>
  );
};
```

#### 2.1.2 型定義ファイルの命名規則

- グローバル型: `src/types/{domain}.ts`
- 機能固有の型: `src/features/{feature}/types.ts`
- APIの型: バックエンドから自動生成された `src/lib/api.ts` を使用

#### 2.1.3 型安全性の徹底

```tsx
// ✅ 良い例: 厳密な型定義
type Status = 'idle' | 'loading' | 'success' | 'error'

// ❌ 悪い例: any型の使用
const handleData = (data: any) => { /* ... */ }

// ✅ 良い例: unknown型を使用して型チェック
const handleData = (data: unknown) => {
  if (isExamData(data)) {
    // 型ガードで安全に使用
  }
}
```

### 2.2. 関数宣言規約

#### 2.2.1 関数の定義方法

**【禁止】function宣言の使用**

```tsx
// ❌ 悪い例: function宣言
function calculateScore(answers: Answer[]): number {
  return answers.filter(a => a.isCorrect).length
}

// ❌ 悪い例: 名前付き関数式
const calculateScore = function(answers: Answer[]): number {
  return answers.filter(a => a.isCorrect).length
}
```

```tsx
// ✅ 良い例: アロー関数を使用
const calculateScore = (answers: Answer[]): number => {
  return answers.filter(a => a.isCorrect).length
}

// ✅ 良い例: 単一式の場合は省略形
const calculateScore = (answers: Answer[]): number =>
  answers.filter(a => a.isCorrect).length
```

#### 2.2.2 コンポーネントの定義

```tsx
// ✅ 良い例: アロー関数でコンポーネントを定義（Propsなし）
export const ExamList = () => {
  return <div>Exam List</div>
}

// ✅ 良い例: propsオブジェクトで受け取る（インライン型定義）
export const ExamCard = (props: {
  exam: { id: string; name: string };
  onClick: (id: string) => void;
}) => {
  return (
    <div onClick={() => props.onClick(props.exam.id)}>
      <h3>{props.exam.name}</h3>
    </div>
  );
};

// ✅ 良い例: 型を別ファイルに定義する場合
import type { ExamListProps } from './types';

export const ExamList = (props: ExamListProps) => {
  return (
    <div>
      {props.exams.map(exam => (
        <button key={exam.id} onClick={() => props.onSelect(exam.id)}>
          {exam.name}
        </button>
      ))}
    </div>
  );
};
```

### 2.3. コンポーネント設計規約

#### 2.3.1 コンポーネントの分類

- **UIコンポーネント** (`src/components/ui/`): shadcn/uiのコンポーネント
- **共通コンポーネント** (`src/components/`): アプリケーション全体で使用
- **機能コンポーネント** (`src/features/{feature}/components/`): 特定機能専用

#### 2.3.2 Propsの命名規則

```tsx
// ✅ 良い例: イベントハンドラーはonXxx形式
type ButtonProps = {
  onClick: () => void
  onHover?: () => void
  isDisabled?: boolean
  children: React.ReactNode
}

// ❌ 悪い例: 一貫性のない命名
type ButtonProps = {
  handleClick: () => void
  hoverCallback?: () => void
  disabled?: boolean
}
```

#### 2.3.3 条件付きレンダリング

```tsx
// ✅ 良い例: 早期リターンを活用
export const ExamDetail = ({ examId }: ExamDetailProps) => {
  const { data: exam, isLoading, error } = useExam(examId)

  if (isLoading) return <LoadingSpinner />
  if (error) return <ErrorMessage error={error} />
  if (!exam) return <NotFound />

  return <div>{exam.name}</div>
}

// ❌ 悪い例: ネストした三項演算子
export const ExamDetail = ({ examId }: ExamDetailProps) => {
  const { data: exam, isLoading, error } = useExam(examId)

  return isLoading ? <LoadingSpinner /> : error ? <ErrorMessage /> : exam ? <div>{exam.name}</div> : <NotFound />
}
```

### 2.4. Hooks規約

#### 2.4.1 カスタムフックの命名

- 必ず `use` で始める
- 具体的な名前を付ける
- ファイル名とフック名を一致させる

```tsx
// ✅ 良い例: src/hooks/useExamList.ts
export const useExamList = () => {
  const { data, error, isLoading } = useSWR('/api/exams', fetcher)
  
  return {
    exams: data,
    isLoading,
    error,
  }
}
```

### 2.5. スタイリング規約

#### 2.5.1 TailwindCSSの使用

```tsx
// ✅ 良い例: クラス名の順序（レイアウト → サイズ → 装飾）
<div className="flex items-center justify-between w-full p-4 bg-white rounded-lg shadow-md">
  <h2 className="text-lg font-bold text-gray-900">Title</h2>
</div>

// ✅ 良い例: 条件付きクラス名はcnヘルパーを使用
import { cn } from '@/lib/utils'

<button className={cn(
  'px-4 py-2 rounded-md',
  isActive && 'bg-blue-500 text-white',
  isDisabled && 'opacity-50 cursor-not-allowed'
)}>
  Submit
</button>
```

### 2.6. ファイル・フォルダ命名規約

- **コンポーネント**: PascalCase (`ExamCard.tsx`)
- **フック**: camelCase、useプレフィックス (`useExamList.ts`)
- **ユーティリティ**: camelCase (`formatDate.ts`)
- **型定義**: camelCase (`exam.ts`, `user.ts`)
- **定数**: UPPER_SNAKE_CASE (`API_ENDPOINTS.ts`)

### 2.7. インポート順序

```tsx
// 1. React関連
import { useState, useEffect } from 'react'

// 2. 外部ライブラリ
import { useSWR } from 'swr'

// 3. 内部モジュール（絶対パス）
import { Button } from '@/components/ui/button'
import { useExamList } from '@/hooks/useExamList'
import type { Exam } from '@/types/exam'

// 4. 相対パス
import { ExamCard } from './ExamCard'
```

## 3. 開発コマンド

```bash
# 開発サーバー起動
pnpm dev

# ビルド
pnpm build

# Lintチェック
pnpm lint

# フォーマット
pnpm format

# テスト実行
pnpm test

# Storybook起動
pnpm storybook
```

## 4. 環境変数

`.env.local`ファイルを作成し、以下の環境変数を設定してください：

```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_FIREBASE_API_KEY=your_api_key
VITE_FIREBASE_AUTH_DOMAIN=your_auth_domain
VITE_FIREBASE_PROJECT_ID=your_project_id
```

## 5. さらなる情報

- [Vite Documentation](https://vite.dev/)
- [React Documentation](https://react.dev/)
- [TailwindCSS Documentation](https://tailwindcss.com/)
- [shadcn/ui Documentation](https://ui.shadcn.com/)
- [Bulletproof React](https://github.com/alan2207/bulletproof-react)
