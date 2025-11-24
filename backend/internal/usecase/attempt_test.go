package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gcp_antigravity/backend/internal/domain"
	"gcp_antigravity/backend/internal/usecase/input"
)

// MockAttemptRepository is a mock implementation of AttemptRepository
type MockAttemptRepository struct {
	mock.Mock
}

func (m *MockAttemptRepository) Save(ctx context.Context, attempt domain.Attempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func (m *MockAttemptRepository) Find(ctx context.Context, id, userID string) (*domain.Attempt, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Attempt), args.Error(1)
}

// MockQuestionRepository is a mock implementation of QuestionRepository
type MockQuestionRepository struct {
	mock.Mock
}

func (m *MockQuestionRepository) BulkCreate(ctx context.Context, questions []domain.Question) error {
	args := m.Called(ctx, questions)
	return args.Error(0)
}

func (m *MockQuestionRepository) FindByExamSetID(ctx context.Context, examSetID string) ([]domain.Question, error) {
	args := m.Called(ctx, examSetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Question), args.Error(1)
}

// MockUserStatsRepository is a mock implementation of UserStatsRepository
type MockUserStatsRepository struct {
	mock.Mock
}

func (m *MockUserStatsRepository) Save(ctx context.Context, stats domain.UserExamStats) error {
	args := m.Called(ctx, stats)
	return args.Error(0)
}

func (m *MockUserStatsRepository) Find(ctx context.Context, userID, examID string) (*domain.UserExamStats, error) {
	args := m.Called(ctx, userID, examID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserExamStats), args.Error(1)
}

// MockTransactionRepository is a mock implementation of TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Run(ctx context.Context, f func(ctx context.Context) error) error {
	// Simply execute the function for testing purposes, ignoring transaction logic
	return f(ctx)
}

func TestUpdateAttempt_MergesAnswers(t *testing.T) {
	// Setup
	mockAttemptRepo := new(MockAttemptRepository)
	mockQuestionRepo := new(MockQuestionRepository)
	mockStatsRepo := new(MockUserStatsRepository)
	mockTxRepo := new(MockTransactionRepository)

	usecase := NewAttemptUsecase(mockQuestionRepo, mockAttemptRepo, mockStatsRepo, mockTxRepo)

	ctx := context.Background()
	userID := "user123"
	attemptID := "attempt123"
	examID := "exam1"
	examSetID := "set1"

	initialAnswers := map[string][]string{
		"q1": {"a"},
		"q2": {"b"},
	}

	existingAttempt := &domain.Attempt{
		ID:             attemptID,
		UserID:         userID,
		ExamID:         examID,
		ExamSetID:      examSetID,
		Status:         domain.StatusInProgress,
		Answers:        initialAnswers,
		TotalQuestions: 10,
		StartedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Mock Find to return the existing attempt
	mockAttemptRepo.On("Find", ctx, attemptID, userID).Return(existingAttempt, nil)

	// Mock Save to capture the updated attempt
	var capturedAttempt domain.Attempt
	mockAttemptRepo.On("Save", ctx, mock.MatchedBy(func(a domain.Attempt) bool {
		capturedAttempt = a
		return true
	})).Return(nil)

	// Execute
	newAnswers := map[string][]string{
		"q2": {"c"}, // Update existing answer
		"q3": {"d"}, // Add new answer
	}
	req := input.UpdateAttemptRequest{
		CurrentIndex: 5,
		Answers:      newAnswers,
	}

	err := usecase.UpdateAttempt(ctx, userID, attemptID, req)

	// Verify
	assert.NoError(t, err)
	mockAttemptRepo.AssertExpectations(t)

	// Check if answers are merged correctly
	assert.Equal(t, 3, len(capturedAttempt.Answers))
	assert.Equal(t, []string{"a"}, capturedAttempt.Answers["q1"]) // Should remain unchanged
	assert.Equal(t, []string{"c"}, capturedAttempt.Answers["q2"]) // Should be updated
	assert.Equal(t, []string{"d"}, capturedAttempt.Answers["q3"]) // Should be added
	assert.Equal(t, 5, capturedAttempt.CurrentIndex)
}

func TestIsCorrect(t *testing.T) {
	tests := []struct {
		name       string
		userAns    []string
		correctAns []string
		expected   bool
	}{
		{
			name:       "同じ要素を持つが順序が異なる場合",
			userAns:    []string{"1", "2", "3"},
			correctAns: []string{"2", "3", "1"},
			expected:   true,
		},
		{
			name:       "完全に一致する場合",
			userAns:    []string{"a", "b"},
			correctAns: []string{"a", "b"},
			expected:   true,
		},
		{
			name:       "要素の数が異なる場合",
			userAns:    []string{"a"},
			correctAns: []string{"a", "b"},
			expected:   false,
		},
		{
			name:       "要素は同じだが重複回数が異なる場合",
			userAns:    []string{"a", "a", "b"},
			correctAns: []string{"a", "b", "b"},
			expected:   false,
		},
		{
			name:       "要素が異なる場合",
			userAns:    []string{"a", "b"},
			correctAns: []string{"c", "d"},
			expected:   false,
		},
		{
			name:       "片方が空のスライス",
			userAns:    []string{},
			correctAns: []string{"a"},
			expected:   false,
		},
		{
			name:       "両方空のスライス",
			userAns:    []string{},
			correctAns: []string{},
			expected:   true,
		},
		{
			name:       "重複がある場合でも順序が異なる場合",
			userAns:    []string{"a", "a", "b"},
			correctAns: []string{"a", "b", "a"},
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isCorrect(tt.userAns, tt.correctAns), tt.name)
		})
	}
}