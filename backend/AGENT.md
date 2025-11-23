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
