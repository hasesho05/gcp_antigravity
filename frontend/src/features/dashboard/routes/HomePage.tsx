import { ExamCard } from "../../exam/components/ExamCard";
import { useExams } from "../api/getExams";
import type { Exam } from "@/types/api";

export function HomePage() {
  const { exams: apiExams, isLoading } = useExams();

  const exams = apiExams?.map((exam: Exam) => ({
    id: exam.id,
    title: exam.name,
    isPro: exam.code.startsWith('P'), // 仮のロジック: Pで始まるコードはProとする
    domainCount: 4, // 仮の値
  })) || [];

  if (isLoading) {
    return <div className="flex justify-center items-center min-h-screen">Loading...</div>;
  }

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      {/* ヘッダー (簡易) */}
      <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/80 backdrop-blur-md">
        <div className="container mx-auto h-16 flex items-center px-4">
          <div className="flex items-center gap-2 font-bold text-xl tracking-tighter">
            <div className="w-6 h-6 rounded bg-linear-to-br from-gcp-blue to-gcp-red flex items-center justify-center text-white text-xs">A</div>
            <span>nearline</span>
          </div>
        </div>
      </header>

      <main className="flex-1 container mx-auto px-4 py-12">
        
        {/* ヒーローセクション */}
        <section className="mb-16 text-center space-y-6 animate-fade-in-up">
          <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight bg-clip-text text-transparent bg-linear-to-r from-primary via-white to-primary/50">
            Master Google Cloud.
          </h1>
          <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
            プロフェッショナルな模擬試験プラットフォーム。
            <br className="hidden md:inline" /> 
            効率的なドメイン別分析で、最短合格を目指しましょう。
          </p>
        </section>

        {/* 試験一覧グリッド */}
        <section>
          <div className="flex items-center justify-between mb-8">
            <h2 className="text-2xl font-bold tracking-tight border-l-4 border-gcp-blue pl-3">
              Available Exams
            </h2>
            <span className="text-sm text-muted-foreground bg-secondary/50 px-3 py-1 rounded-full">
              {exams.length} Exams
            </span>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {exams.map((exam, index) => (
              <ExamCard 
                key={exam.id} 
                exam={exam} 
                // スクロールに合わせてふわっと表示させる遅延アニメーション
                className="animate-fade-in-up opacity-0 fill-mode-forwards"
                style={{ animationDelay: `${index * 80}ms` }}
              />
            ))}
          </div>
        </section>

      </main>
    </div>
  );
}
