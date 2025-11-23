package domain

import "time"

// Question は1つの問題を表すマスターデータです。
type Question struct {
	ID                 string         `json:"id" firestore:"id"`                                   // Document ID (e.g., "PCD_SET1_001")
	ExamID             string         `json:"examId" firestore:"exam_id"`                          // 資格ID (e.g., "professional_cloud_developer")
	ExamSetID          string         `json:"examSetId" firestore:"exam_set_id"`                   // 模擬試験セットID (e.g., "practice_exam_1")
	ExamCode           string         `json:"examCode" firestore:"exam_code"`                      // 資格コード (e.g., "PCD")
	QuestionText       string         `json:"question" firestore:"question_text"`                  // HTML string
	QuestionType       string         `json:"questionType" firestore:"question_type"`              // "multiple-choice" or "multi-select"
	Options            []AnswerOption `json:"answerOptions" firestore:"options"`                   // 選択肢リスト
	CorrectAnswers     []string       `json:"correctAnswers" firestore:"correct_answers"`          // 正解のOption IDリスト
	OverallExplanation string         `json:"overallExplanation" firestore:"overall_explanation"`  // 全体の解説 (HTML)
	Domain             string         `json:"domain" firestore:"domain"`                           // 分野 (e.g. "Compute")
	ImageURL           string         `json:"imageUrl,omitempty" firestore:"image_url,omitempty"`  // 解説図などのURL
	ReferenceURLs      []string       `json:"referenceUrls,omitempty" firestore:"reference_urls"`  // 参考リンク
	CreatedAt          time.Time      `json:"createdAt" firestore:"created_at"`
}

// AnswerOption は問題の個々の選択肢です。
type AnswerOption struct {
	ID          string `json:"id" firestore:"id"`                   // "a", "b", "c", "d" or UUID
	Text        string `json:"answer" firestore:"text"`             // 選択肢の文言
	Explanation string `json:"explanation" firestore:"explanation"` // この選択肢ごとの解説
}
