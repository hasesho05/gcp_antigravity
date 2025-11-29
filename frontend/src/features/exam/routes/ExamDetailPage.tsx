import { useParams } from 'react-router-dom';
import { useExam } from '../api/getExam';
import { useExamSets } from '../api/getExamSets';
import { useExamSetStats } from '../api/getExamSetStats';
import { Progress } from '@/components/ui/progress';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Loader2, PlayCircle, BookOpen, BarChart2 } from 'lucide-react';

export const ExamDetailPage = () => {
  const { examId } = useParams<{ examId: string }>();
  const { exam, isLoading: isExamLoading } = useExam(examId);
  const { examSets, isLoading: isSetsLoading } = useExamSets(examId);
  const { stats } = useExamSetStats(examId);

  if (isExamLoading || isSetsLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (!exam) {
    return <div className="p-8 text-center">Exam not found</div>;
  }

  const getSetStats = (setId: string) => stats?.find((s) => s.examSetId === setId);

  return (
    <div className="min-h-screen bg-background pb-20">
      {/* Top Section: Question Answer Placeholder */}
      <div className="bg-muted/30 border-b border-border/50">
        <div className="container mx-auto px-4 py-12">
          <div className="flex flex-col items-center justify-center text-center space-y-4">
             <div className="p-6 rounded-2xl bg-card border border-border/50 shadow-sm max-w-2xl w-full">
                <p className="text-muted-foreground">
                  以下の模擬試験セットを選択して、ここで問題演習を開始してください。
                </p>
                <div className="mt-4 h-32 bg-muted/50 rounded-lg flex items-center justify-center border border-dashed border-border">
                    <span className="text-sm text-muted-foreground">問題演習インターフェース（準備中）</span>
                </div>
             </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 -mt-8 relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-8">
             {/* Exam Info Card */}
            <Card className="border-border/50 shadow-lg bg-card/80 backdrop-blur-sm">
              <CardHeader>
                <div className="flex items-start justify-between">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight mb-2">{exam.name}</h1>
                        <p className="text-muted-foreground">{exam.description}</p>
                    </div>
                    {exam.imageUrl && (
                        <img src={exam.imageUrl} alt={exam.name} className="w-16 h-16 object-contain" />
                    )}
                </div>
                
                <div className="flex items-center gap-4 mt-6 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                        <BarChart2 className="w-4 h-4" />
                        <span>プロフェッショナルレベル</span>
                    </div>
                    <div className="flex items-center gap-1">
                        <BookOpen className="w-4 h-4" />
                        <span>{examSets?.length || 0} 模擬試験セット</span>
                    </div>
                </div>
              </CardHeader>
            </Card>

            {/* Tabs Section */}
            <Tabs defaultValue="content" className="w-full">
              <TabsList className="grid w-full grid-cols-3 lg:w-[400px]">
                <TabsTrigger value="content">コース内容</TabsTrigger>
                <TabsTrigger value="overview">概要</TabsTrigger>
                <TabsTrigger value="reviews">レビュー</TabsTrigger>
              </TabsList>
              <TabsContent value="content" className="mt-6 space-y-4">
                {examSets?.map((set) => {
                  const setStats = getSetStats(set.id);
                  const isStarted = setStats && setStats.totalAttempts > 0;
                  const latest = setStats?.latestAttempt;
                  const isInProgress = latest?.status === 'in_progress';

                  return (
                  <Card key={set.id} className="group hover:border-primary/50 transition-colors cursor-pointer">
                    <CardContent className="p-6 flex items-center justify-between">
                      <div className="flex items-center gap-4 flex-1">
                        <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-colors shrink-0">
                            <PlayCircle className="w-5 h-5" />
                        </div>
                        <div className="flex-1">
                            <h3 className="font-semibold text-lg">{set.name}</h3>
                            <p className="text-sm text-muted-foreground mb-2">{set.description}</p>
                            
                            {isStarted && (
                              <div className="flex items-center gap-4 text-xs text-muted-foreground mt-2">
                                <span>試行回数: {setStats.totalAttempts}回</span>
                                {isInProgress && (
                                  <div className="flex items-center gap-2 flex-1 max-w-[200px]">
                                    <Progress value={latest.progress} className="h-2" />
                                    <span>{latest.progress}%</span>
                                  </div>
                                )}
                                {latest?.status === 'completed' && (
                                  <span className="text-green-500 font-medium">完了 (スコア: {latest.score}%)</span>
                                )}
                              </div>
                            )}
                        </div>
                      </div>
                      <Button variant={isInProgress ? "default" : "ghost"} size="sm">
                        {isInProgress ? "再開" : "開始"}
                      </Button>
                    </CardContent>
                  </Card>
                )})}
                {(!examSets || examSets.length === 0) && (
                    <div className="text-center py-12 text-muted-foreground">
                        模擬試験セットはまだありません。
                    </div>
                )}
              </TabsContent>
              <TabsContent value="overview">
                <Card>
                  <CardContent className="p-6">
                    <h3 className="text-lg font-semibold mb-2">この試験について</h3>
                    <p className="text-muted-foreground">
                      この試験は、ビジネス目標を達成するための堅牢で安全、スケーラブル、高可用性、動的なソリューションを設計、開発、管理する能力を評価します。
                    </p>
                  </CardContent>
                </Card>
              </TabsContent>
              <TabsContent value="reviews">
                <Card>
                  <CardContent className="p-6 text-center text-muted-foreground">
                    レビューは準備中です。
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>

          {/* Sidebar */}
          <div className="lg:col-span-1">
            <Card className="sticky top-24 border-border/50 shadow-lg">
                <CardContent className="p-6 space-y-6">
                    <div>
                        <h3 className="font-semibold mb-2">学習を開始しますか？</h3>
                        <p className="text-sm text-muted-foreground mb-4">
                            進捗を記録し、本番環境に近い形式でシミュレーションを行います。
                        </p>
                        <Button className="w-full btn-gradient text-white" size="lg">
                            演習を開始する
                        </Button>
                    </div>
                    
                    <div className="pt-6 border-t border-border">
                        <h4 className="text-sm font-semibold mb-3">含まれる内容:</h4>
                        <ul className="space-y-2 text-sm text-muted-foreground">
                            <li className="flex items-center gap-2">
                                <span className="w-1.5 h-1.5 rounded-full bg-green-500" />
                                {examSets?.length || 0} 模擬試験セット
                            </li>
                            <li className="flex items-center gap-2">
                                <span className="w-1.5 h-1.5 rounded-full bg-green-500" />
                                詳細な解説
                            </li>
                            <li className="flex items-center gap-2">
                                <span className="w-1.5 h-1.5 rounded-full bg-green-500" />
                                パフォーマンス追跡
                            </li>
                        </ul>
                    </div>
                </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
};
