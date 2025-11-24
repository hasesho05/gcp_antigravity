package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gcp_antigravity/backend/internal/domain"
)

func main() {
	// Create dummy data for all domain models
	now := time.Now()
	
	user := domain.User{
		ID:                 "user123",
		Email:              "test@example.com",
		Role:               domain.RolePro,
		SubscriptionStatus: domain.SubActive,
		CreatedAt:          now,
	}

	question := domain.Question{
		ID:           "PCD_SET1_001",
		ExamID:       "professional_cloud_developer",
		ExamSetID:    "practice_exam_1",
		ExamCode:     "PCD",
		QuestionText: "<p>What is Cloud Run?</p>",
		QuestionType: "multiple-choice",
		Options: []domain.AnswerOption{
			{ID: "a", Text: "Serverless container platform", Explanation: "Correct."},
			{ID: "b", Text: "VM", Explanation: "Incorrect."},
		},
		CorrectAnswers:     []string{"a"},
		OverallExplanation: "<p>Cloud Run is managed serverless...</p>",
		Domain:             "Compute",
		CreatedAt:          now,
	}

	attempt := domain.Attempt{
		ID:             "attempt123",
		UserID:         "user123",
		ExamID:         "professional_cloud_developer",
		ExamSetID:      "practice_exam_1",
		Status:         domain.StatusInProgress,
		Score:          0,
		TotalQuestions: 50,
		CurrentIndex:   5,
		Answers:        map[string][]string{"PCD_SET1_001": {"a"}},
		StartedAt:      now,
		UpdatedAt:      now,
	}

	stats := domain.UserExamStats{
		ExamID:                 "professional_cloud_developer",
		UserID:                 "user123",
		TotalAttempts:          1,
		TotalScore:             85,
		TotalQuestionsAnswered: 100,
		DomainStats: map[string]domain.DomainScore{
			"Compute": {DomainName: "Compute", CorrectCount: 8, TotalCount: 10, AccuracyRate: 80},
		},
		LastTakenAt: now,
	}

	// Combine into a single object for generation
	output := map[string]interface{}{
		"User":          user,
		"Question":      question,
		"Attempt":       attempt,
		"UserExamStats": stats,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}
