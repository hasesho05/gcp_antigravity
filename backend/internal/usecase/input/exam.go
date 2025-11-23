package input

type UploadQuestionsRequest struct {
	ExamID    string         `json:"examId"`
	ExamSetID string         `json:"examSetId"`
	ExamCode  string         `json:"examCode"`
	Questions []QuestionInput `json:"questions"`
}

type QuestionInput struct {
	Index              int      `json:"index"` // 1-based index for ID generation
	QuestionText       string   `json:"question"`
	QuestionType       string   `json:"questionType"`
	Options            []OptionInput `json:"options"`
	CorrectAnswers     []string `json:"correctAnswers"`
	OverallExplanation string   `json:"overallExplanation"`
	Domain             string   `json:"domain"`
	ImageURL           string   `json:"imageUrl"`
	ReferenceURLs      []string `json:"referenceUrls"`
}

type OptionInput struct {
	ID          string `json:"id"`
	Text        string `json:"answer"`
	Explanation string `json:"explanation"`
}
