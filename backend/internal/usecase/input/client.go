package input

type CreateAttemptRequest struct {
	ExamID    string `json:"examId"`
	ExamSetID string `json:"examSetId"`
}

type UpdateAttemptRequest struct {
	CurrentIndex int                 `json:"currentIndex"`
	Answers      map[string][]string `json:"answers"` // Key: QuestionID
}

type CompleteAttemptRequest struct {
	Answers map[string][]string `json:"answers"` // Final answers
}
