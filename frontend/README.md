# GCP認定資格模擬試験プラットフォーム - Frontend

このプロジェクトは、GCP認定資格模擬試験プラットフォームのフロントエンドアプリケーションです。Bulletproof Reactアーキテクチャパターンに基づいて設計されています。

## Tech Stack

- **Framework**: React 18+ with Vite
- **Language**: TypeScript
- **Styling**: TailwindCSS + shadcn/ui
- **Testing**: Vitest + Testing Library
- **Storybook**: UIコンポーネントカタログ
- **Data Fetching**: useSWR
- **State Management**: Zustand
- **Routing**: React Router v6
- **Linter/Formatter**: Biome

## Getting Started

### Prerequisites

- Node.js 18+
- pnpm 8+

### Installation

```bash
pnpm install
```

### Environment Setup

環境変数を設定します:

```bash
cp .env.example .env
```

`.env`ファイルを編集して、必要な環境変数を設定してください。

### Available Scripts

#### 開発

```bash
pnpm dev
```

開発サーバーを起動します。デフォルトで `http://localhost:5173` で起動します。

#### ビルド

```bash
pnpm build
```

本番用ビルドを作成します。

#### プレビュー

```bash
pnpm preview
```

ビルドされたアプリケーションをローカルでプレビューします。

#### テスト

```bash
pnpm test          # テストを実行
pnpm test:ui       # Vitest UIでテストを実行
```

#### Linting & Formatting

```bash
pnpm lint          # コードをチェック
pnpm lint:fix      # 自動修正可能な問題を修正
pnpm format        # コードをフォーマット
```

#### Storybook

```bash
pnpm storybook           # Storybookを起動
pnpm build-storybook     # Storybookをビルド
```

#### API型生成

```bash
pnpm gen:types
```

バックエンドのJSON定義からTypeScript型を生成します。

## Project Structure

```
src/
├── assets/              # 静的アセット
├── components/          # 共有UIコンポーネント
│   ├── ui/              # shadcn/uiコンポーネント
│   ├── elements/        # 共有原子コンポーネント
│   └── layouts/         # ページレイアウト
├── config/              # アプリケーション設定
├── features/            # 機能モジュール
│   ├── auth/            # 認証機能
│   ├── exam/            # 試験機能
│   ├── dashboard/       # ダッシュボード
│   └── admin/           # 管理者機能
├── hooks/               # 汎用カスタムフック
├── lib/                 # ライブラリ設定・ラッパー
├── providers/           # アプリケーションプロバイダー
├── routes/              # ルーティング設定
├── stores/              # グローバル状態管理
├── test/                # テストセットアップ
├── types/               # 型定義（自動生成）
└── utils/               # ユーティリティ関数
```

## Development Guidelines

### Adding a New Feature

1. `src/features/`に新しいディレクトリを作成
2. 必要な型定義が`src/types/api.ts`にあるか確認
3. `api/`でuseSWRラッパーフックを作成
4. `components/`でUIコンポーネントを作成
5. `routes/`でページコンポーネントを作成
6. `src/routes/index.tsx`にルートを追加

### Component Development

- Presentational componentsは`features/{feature}/components/`に配置
- 共有コンポーネントは`src/components/`に配置
- Storybookでコンポーネントを開発・カタログ化

### API Integration

- `useEffect`でのデータフェッチは禁止
- 必ず`useSWR`を使用
- API型は`src/types/api.ts`から自動生成

### Type Synchronization

バックエンドの型定義が更新された場合:

1. Backend: `make generate-sample`
2. Frontend: `pnpm gen:types`

## Documentation

詳細な設計ドキュメントは`PROJECT_DESIGN.md`を参照してください。
