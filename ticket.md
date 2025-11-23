ID,優先度,担当レイヤー,タスク名,詳細/完了条件
I-101,High,[SETUP],Goプロジェクト初期化と依存関係の追加,go.mod の作成、cockroachdb/errors、Firebase Admin SDK for Goなどの必須ライブラリをインストール。 (Done)
I-102,High,[SETUP],HTTPサーバーとCORS設定,cmd/api/main.go でのHTTPサーバー（net/http）のセットアップと、Cloud Run運用を見据えたCORS設定を実装。 (Done)
I-103,High,[INFRA],Firestoreクライアントの共通実装,infra/firestore/client.go にFirestoreクライアントの初期化、接続、エラーハンドリングの共通処理を実装。 (Done)
I-104,High,[SETUP],DIコンテナのセットアップ,"main.goでHandler, Usecase, RepositoryImpl間の依存性注入（DI）を実装し、構造を確立する。"
I-105,High,[DOMAIN],ドメインモデルの定義,"internal/domain/{question, attempt, stats}.go にGo構造体を定義し、json タグと firestore タグを全て正確に記述する。" (Done)
I-106,High,[SETUP],Quicktype連携用Makefile実装,Makefile に generate-sample ターゲットを実装。scripts/dump_json.go を作成し、定義したDomain/DTOの構造体からJSONサンプルを出力できるようにする。 (Done)
A-201,High,[REPO],ExamRepositoryインターフェース定義,internal/repository/repository.go に BulkCreateQuestions などAdminに必要なメソッドを定義。 (Done)
A-202,High,[IMPL],BulkCreateQuestions の実装,repository_impl/exam.go に問題一括登録処理を実装。FirestoreのBatch Writeを利用し、コストを抑える。 (Done)
A-203,High,[USECASE],問題入稿ロジックの実装,usecase/exam.go に UploadQuestions メソッドを実装。Adminからの入力JSON（DTO）をDomainモデルに変換し、バリデーションとID生成を行う。 (Done)
A-204,High,[HANDLER],Admin問題入稿APIの実装,handler/admin/handler.go に POST /admin/exams/{examID}/sets/{setID}/questions エンドポイントを実装。 (Done)
C-301,Mid,[REPO],Client関連インターフェース定義,repository/repository.go に GetQuestionsByExamSetID、SaveAttempt、GetAttempt、UpdateStats を追加定義。 (Done)
C-302,Mid,[IMPL],問題取得/Attempt保存の実装,repository_impl/exam.go に上記C-301で定義した基本的なCRUD操作を実装する。 (Done)
C-303,Mid,[USECASE],問題取得とAttempt開始ロジック,usecase/exam.go に GetExamQuestions と StartAttempt を実装。 (Done)
C-304,Mid,[HANDLER],問題一覧/Attempt開始APIの実装,handler/client/handler.go に GET /exams/.../questions と POST /users/me/attempts エンドポイントを実装。 (Done)
C-305,Low,[USECASE],採点とStats更新ロジックの実装,usecase/exam.go に CompleteAttempt の核となる採点ロジックと、**Stats更新（Firestore Transaction利用）**ロジックを実装。
C-306,Low,[HANDLER],進捗保存/試験完了APIの実装,PUT /users/me/attempts/{attemptID} と POST /users/me/attempts/{attemptID}/complete を実装。
C-307,Low,[UX],苦手分野集計結果の取得API,ユーザーのダッシュボード表示用に GET /users/me/stats/{examID} を実装する。