package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"nearline/backend/internal/domain"
	"nearline/backend/internal/infra/firestore"

	"github.com/joho/godotenv"
	"github.com/samber/lo"
)

// Configuration
// 実行前にここを変更してください
const (
	InputFile       = "cmd/seed_questions/source.json"
	TargetExamID    = "professional-cloud-developer" // FirestoreのExam ID
	TargetExamCode  = "PCD"                          // 問題IDのプレフィックスに使用
	QuestionsPerSet = 50
)

// Domain Mapping (Small Category -> Large Category)
// 試験ごとに手動で変更してください
var DomainMapping = map[string]string{
	"IAMとセキュリティ":     "Identity and Security",
	"モニタリングと運用管理": "Monitoring and Operations",
	// 必要に応じて追加
}

// Input Structs
type InputJSON struct {
	Questions []InputQuestion `json:"questions"`
}

type InputQuestion struct {
	Question           string              `json:"question"`
	QuestionType       string              `json:"questionType"`
	AnswerOptions      []InputAnswerOption `json:"answerOptions"`
	OverallExplanation string              `json:"overallExplanation"`
	CorrectAnswers     string              `json:"correctAnswers"` // "1" or "1,3"
	Domain             string              `json:"domain"`
}

type InputAnswerOption struct {
	Answer      string `json:"answer"`
	Explanation string `json:"explanation"`
}

func main() {
	// 1. .env の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println(".env ファイルが見つかりません。環境変数に依存します。")
	}

	// 2. Firestore への接続
	ctx := context.Background()
	client := firestore.NewClient(ctx)
	defer client.Close()

	// 3. JSON の読み込み
	data, err := os.ReadFile(InputFile)
	if err != nil {
		log.Fatalf("入力ファイルの読み込みに失敗しました: %v", err)
	}

	var input InputJSON
	if err := json.Unmarshal(data, &input); err != nil {
		log.Fatalf("JSON のパースに失敗しました: %v", err)
	}

	log.Printf("%s から %d 問の問題を読み込みました", InputFile, len(input.Questions))

	// 4. 問題のバリデーションとマッピング
	var validQuestions []*domain.Question
	
	for i, q := range input.Questions {
		// ドメインのマッピング
		largeDomain, ok := DomainMapping[q.Domain]
		if !ok {
			log.Printf("警告: ドメイン '%s' がマッピングに見つかりません。元のドメインを使用します。", q.Domain)
			largeDomain = q.Domain
		}

		// 選択肢と正解のマッピング
		options := make([]domain.AnswerOption, 0, len(q.AnswerOptions))
		for j, opt := range q.AnswerOptions {
			options = append(options, domain.AnswerOption{
				ID:          fmt.Sprintf("%d", j+1), // "1", "2", "3"...
				Text:        opt.Answer,
				Explanation: opt.Explanation,
			})
		}

		// 正解のパース
		correctAnswersStr := strings.Split(q.CorrectAnswers, ",")
		correctAnswers := lo.Map(correctAnswersStr, func(s string, _ int) string {
			return strings.TrimSpace(s)
		})

		// 正解のバリデーション
		hasError := false

		// 4-1. 問題タイプに基づくチェック
		if q.QuestionType == "multiple-choice" && len(correctAnswers) > 1 {
			log.Printf("エラー: 問題 %d は multiple-choice ですが、正解が複数あります: %s", i, q.CorrectAnswers)
			hasError = true
		}

		// 4-2. 数値チェックと範囲チェック
		for _, ans := range correctAnswers {
			num, err := strconv.Atoi(ans)
			if err != nil {
				log.Printf("エラー: 問題 %d の正解 '%s' が数値ではありません。", i, ans, i)
				hasError = true
				break
			}
			
			if num < 1 || num > len(options) {
				log.Printf("エラー: 問題 %d の正解 '%s' が選択肢の範囲(1-%d)外です。", i, ans, len(options))
				hasError = true
				break
			}
		}

		if hasError {
			log.Printf("問題 %d をスキップします。", i)
			continue
		}

		// ドメイン質問の作成
		// ID と SetID は後で設定するため、一時的にプレースホルダーを使用
		tempID := fmt.Sprintf("TEMP_%d", i)
		tempSetID := "TEMP_SET"
		
		domainQ, err := domain.NewQuestion(
			tempID,
			TargetExamID,
			tempSetID,
			TargetExamCode,
			q.Question,
			q.QuestionType,
			q.OverallExplanation,
			largeDomain,
			"", // ImageURL
			options,
			correctAnswers,
			nil, // ReferenceURLs
			time.Now(),
		)
		if err != nil {
			log.Printf("エラー: 問題 %d の作成に失敗しました: %v", i, err)
			continue
		}

		validQuestions = append(validQuestions, domainQ)
	}

	// 5. 均等な分散 (ドメインごとのグループ化 -> インターリーブ)
	// ドメインごとに問題をグループ化
	questionsByDomain := lo.GroupBy(validQuestions, func(q *domain.Question) string {
		return q.Domain
	})

	// インターリーブ用のスライスのスライスを作成
	var domainGroups [][]*domain.Question
	for _, qs := range questionsByDomain {
		// 各ドメイングループ内でシャッフル
		domainGroups = append(domainGroups, lo.Shuffle(qs))
	}

	// インターリーブしてドメインを均等に混ぜる
	balancedQuestions := lo.Interleave(domainGroups...)
	log.Printf("%d 問の問題を %d のドメインにわたってバランス良く配置しました", len(balancedQuestions), len(domainGroups))

	// 6. セットへの分割
	chunks := lo.Chunk(balancedQuestions, QuestionsPerSet)
	log.Printf("%d セットの問題を作成しました", len(chunks))

	// 7. Firestore への保存 (BulkWriter)
	bulkWriter := client.BulkWriter(ctx)

	for setIndex, chunk := range chunks {
		// Set ID の生成: practice_exam_{n} 形式
		setID := fmt.Sprintf("practice_exam_%d", setIndex+1) 
		log.Printf("セット %d (ID: %s) の %d 問を処理中...", setIndex+1, setID, len(chunk))

		// ExamSet ドキュメントの作成
		examSet := domain.ExamSet{
			ID:          setID,
			ExamID:      TargetExamID,
			Name:        fmt.Sprintf("Practice Exam %d", setIndex+1),
			Description: fmt.Sprintf("%d questions covering all domains", len(chunk)),
			QuestionIDs: []string{}, // 後で埋めることも可能だが、ここでは省略
			CreatedAt:   time.Now(),
		}
		
		// exams/{examID}/sets/{setID} に保存
		setRef := client.Collection("exams").Doc(TargetExamID).Collection("sets").Doc(setID)
		_, err := bulkWriter.Set(setRef, examSet)
		if err != nil {
			log.Printf("エラー: ExamSet の保存に失敗しました (ID: %s): %v", setID, err)
		}

		for qIndex, q := range chunk {
			// 問題 ID の生成: ExamCode_SetIndex_QuestionIndex (例: PCD_SET1_001)
			q.ID = fmt.Sprintf("%s_SET%d_%03d", TargetExamCode, setIndex+1, qIndex+1)
			q.ExamSetID = setID

			// exams/{examID}/sets/{setID}/questions/{questionID} に保存
			docRef := setRef.Collection("questions").Doc(q.ID)
			_, err := bulkWriter.Set(docRef, q)
			if err != nil {
				log.Printf("エラー: BulkWriter への追加に失敗しました (ID: %s): %v", q.ID, err)
			}
		}
	}

	// BulkWriter のフラッシュとクローズ
	bulkWriter.Flush()
	
	log.Println("シーディングが正常に完了しました！")
}
