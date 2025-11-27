package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"nearline/backend/internal/domain"
	"nearline/backend/internal/infra/firestore"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	ctx := context.Background()
	client := firestore.NewClient(ctx)
	defer client.Close()

	exams := []domain.Exam{
		{
			ID:          "cloud-digital-leader",
			Code:        "CDL",
			Name:        "Cloud Digital Leader",
			Description: "Google Cloud のコアプロダクトとサービスに関する知識、およびそれらが組織にどのように利益をもたらすかを理解していることを示します。",
			ImageURL:    "/images/exams/cdl.png",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "associate-cloud-engineer",
			Code:        "ACE",
			Name:        "Associate Cloud Engineer",
			Description: "アプリケーションのデプロイ、オペレーションのモニタリング、エンタープライズ ソリューションの管理を行います。",
			ImageURL:    "/images/exams/ace.png",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "professional-cloud-architect",
			Code:        "PCA",
			Name:        "Professional Cloud Architect",
			Description: "Google Cloud 技術を活用した、安全でスケーラブル、かつ可用性の高い堅牢なソリューションを設計、開発、管理する能力を評価します。",
			ImageURL:    "/images/exams/pca.png",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "professional-cloud-developer",
			Code:        "PCD",
			Name:        "Professional Cloud Developer",
			Description: "スケーラブルで可用性の高いアプリケーションを構築、デプロイ、管理する能力を評価します。",
			ImageURL:    "/images/exams/pcd.png",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "professional-data-engineer",
			Code:        "PDE",
			Name:        "Professional Data Engineer",
			Description: "データの収集、変換、公開によって、データ主導の意思決定を可能にします。",
			ImageURL:    "/images/exams/pde.png",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "professional-cloud-devops-engineer",
			Code:        "PDOE",
			Name:        "Professional Cloud DevOps Engineer",
			Description: "効率的な開発運用パイプラインを構築し、サービスの信頼性を維持する能力を評価します。",
			ImageURL:    "/images/exams/pdoe.png",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "professional-cloud-security-engineer",
			Code:        "PCSE",
			Name:        "Professional Cloud Security Engineer",
			Description: "Google Cloud 上で安全なインフラストラクチャを設計、実装する能力を評価します。",
			ImageURL:    "/images/exams/pcse.png",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "professional-cloud-network-engineer",
			Code:        "PCNE",
			Name:        "Professional Cloud Network Engineer",
			Description: "Google Cloud 上でネットワーク アーキテクチャを実装、管理する能力を評価します。",
			ImageURL:    "/images/exams/pcne.png",
			CreatedAt:   time.Now(),
		},
		{
			ID:          "professional-machine-learning-engineer",
			Code:        "PMLE",
			Name:        "Professional Machine Learning Engineer",
			Description: "ML モデルの構築、評価、本番環境へのデプロイ、および最適化を行う能力を評価します。",
			ImageURL:    "/images/exams/pmle.png",
			CreatedAt:   time.Now(),
		},
	}

	for _, exam := range exams {
		_, err := client.Collection("exams").Doc(exam.ID).Set(ctx, exam)
		if err != nil {
			log.Printf("Failed to seed exam %s: %v\n", exam.Name, err)
		} else {
			fmt.Printf("Seeded exam: %s\n", exam.Name)
		}
	}

	fmt.Println("Seeding completed.")
}
