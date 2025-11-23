package input

type CreateAttemptRequest struct {
	ExamID    string `json:"examId"`
	ExamSetID string `json:"examSetId"`
}
