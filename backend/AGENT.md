# Go Backend Development Guidelines (for Agent)

This document outlines key development principles and patterns to follow for this project, based on past code reviews. The agent should adhere to these guidelines when modifying the codebase.

## 1. Architecture & Design

- **Single Responsibility Principle (SRP) for Repositories**: Each repository should be responsible for managing a single domain entity (aggregate root). Avoid creating large, monolithic repositories that handle multiple entities.
    - **Good**: `QuestionRepository`, `AttemptRepository`, `UserStatsRepository`
    - **Bad**: A single `ExamRepository` handling questions, attempts, and stats.

- **Transaction Management in Usecase Layer**: Business transactions that involve multiple repositories should be managed in the usecase layer, not the repository layer. The usecase layer is responsible for orchestrating the steps of a business operation, including beginning and ending a transaction.

## 2. Robustness & Error Handling

- **Use UUIDs for Unique IDs**: When generating unique identifiers for entities (like `AttemptID`), use a standard UUID library (`github.com/google/uuid`) instead of timestamp-based or composite keys. This guarantees uniqueness and avoids potential collisions.

- **Graceful Initialization**: Critical components like database clients or authentication clients should not cause a `panic` on initialization failure. Instead, return a wrapped error from the initialization function (e.g., `run()`) and allow the `main` function to log the fatal error and exit gracefully.

- **Correct Firestore Error Handling**:
    - To check if a document was not found, use `status.Code(err) == codes.NotFound` from `google.golang.org/grpc/status` and `google.golang.org/grpc/codes`.
    - Do not rely on checking `doc.Exists()` after an error has already occurred, as the `doc` object may be unreliable.

## 3. API & Data Handling

- **Consistent API Responses**:
    - For `GET` requests that return a resource, if the resource does not exist, return a `200 OK` with a default or empty state of the object (e.g., an empty `UserExamStats` object) rather than a `404 Not Found`. This simplifies frontend logic.
    - Do not return hardcoded JSON strings. Use `json.NewEncoder` with a `struct` or `map` to ensure well-formed JSON responses.

## 4. Efficiency

- **Avoid Redundant Data Fetching**: Fetch data only when necessary. If a value (like the total number of questions in an exam) is needed in a later step, fetch it at the beginning of the process (e.g., when an `Attempt` is created) and store it within the relevant domain object. Do not re-fetch the same data in a later function call.

## 5. Domain Object Creation

- **Use Constructors for Domain Objects**: Always use a constructor function (e.g., `domain.NewQuestion()`) defined in the domain layer to create new instances of domain objects. This ensures that objects are always created in a valid state.
- **Validation in Constructors**: Constructors should validate their arguments and return an error if any validation fails. This prevents the creation of invalid domain objects.
- **Avoid Direct Struct Initialization**: Do not initialize domain structs directly from other layers (e.g., `&domain.Question{...}`). This bypasses validation and can lead to inconsistent or invalid object states.
- **コンストラクタの使用の徹底**: ドメインオブジェクトを初期化する際は、構造体を直接初期化する (`&domain.UserExamStats{...}`) のではなく、必ずコンストラクタ（例: `domain.NewUserExamStats()`）を使用してください。これにより、オブジェクトの一貫性とドメインルールが保証されます。

## 6. Utility Functions

- **Pointer Helpers**: Use `util.ToPointer` and `util.FromPointer` for converting between values and pointers, especially when dealing with optional fields or interfacing with external libraries that require pointers.
- **`lo.Map` Helper**: When using `github.com/samber/lo`'s `Map` function and the index is not required in the iteratee, prefer using `util.Map` for cleaner code.

## 7. Usecase Layer Refactoring

- **Use Input/Output DTOs**: Usecase methods should not take domain objects as arguments or return them directly. Instead, use dedicated Data Transfer Objects (DTOs) defined in the `usecase/input` and `usecase/output` packages.
- **Clear Separation**: This practice creates a clear separation between the application's core business logic (domain) and its orchestration layer (usecase). The `input` objects encapsulate the parameters required for a usecase, while the `output` objects format the data for the presentation layer (e.g., handlers).
- **Validation in Input Constructors**: Validation of parameters should be performed within the constructor of the `input` DTO (e.g., `input.NewCompleteAttempt()`). This ensures that the usecase always receives valid data.
- **Example**:
    - **Before**: `func (u *myUsecase) DoSomething(ctx context.Context, userID string, param1 int) (*domain.MyObject, error)`
    - **After**: `func (u *myUsecase) DoSomething(ctx context.Context, input *input.DoSomething) (*output.MyObject, error)`

## 8. Domain Models & Data Structure Quick Reference

### 🧩 Domain Models

アプリケーションの中核となるドメインモデルの解説です。

#### 1. Exam (資格試験)
GCP認定試験そのものを表します（例: "Professional Cloud Developer"）。
- `ID`: 資格ID (e.g. `professional-cloud-developer`)
- `Code`: 資格コード (e.g. `PCD`)

#### 2. ExamSet (模擬試験セット)
1つの資格試験に含まれる、50問1セットの模擬試験単位です。
- `ID`: セットID (e.g. `practice_exam_1`)
- `ExamID`: 親となる資格ID

#### 3. Question (問題)
個々の問題データです。
- `ID`: 問題ID (e.g. `PCD_SET1_001`)
- `QuestionType`: `multiple-choice` (単一選択) または `multi-select` (複数選択)
- `Domain`: 出題分野 (e.g. "Identity and Security")
- `OverallExplanation`: 全体の解説

#### 4. User (ユーザー)
アプリケーションの利用者です。Firebase Authと連携します。
- `Role`: `free` (無料), `pro` (有料), `admin` (管理者)
- `SubscriptionStatus`: サブスクリプションの状態

#### 5. Attempt (受験)
ユーザーが模擬試験を1回受験した履歴を表します。
- `Status`: `in_progress` (受験中), `paused` (中断), `completed` (完了)
- `Answers`: ユーザーの回答状況 (Map形式)
- `Score`: 正解数

#### 6. UserExamStats (成績統計)
ユーザーの資格ごとの累積成績です。
- `DomainStats`: 分野ごとの正答率などの統計情報
- `AccuracyRate`: 全体の正答率

### 🔥 Firestore Data Structure

```
exams/{examID}
  ├── sets/{setID} (ExamSet)
  │     ├── questions/{questionID} (Question)
```

- **ExamSet**: 模擬試験のセット（例: "Practice Exam 1"）
- **Question**: 個々の問題データ