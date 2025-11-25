import { Card, CardContent, CardFooter, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Lock, Play, BarChart3 } from "lucide-react";

type Exam = {
  id: string;
  title: string;
  isPro: boolean;
  domainCount: number; // 例: 4つの分野
};

// 実際のロゴ画像のパス (publicフォルダかassetsに配置)
const LOGO_SRC = "/gcp_pca_logo_56c37d60ee.webp"; 

export function ExamCard({ exam, className, ...props }: { exam: Exam } & React.HTMLAttributes<HTMLDivElement>) {
  return (
    <Card className={`group relative overflow-hidden border-border/50 bg-card hover:bg-accent/40 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5 hover:-translate-y-1 ${className}`} {...props}>
      {/* 資格ロゴ部分 */}
      <CardHeader className="flex flex-row items-center gap-4 pb-2">
        <div className="relative h-14 w-14 shrink-0 overflow-hidden rounded-full border border-border/50 bg-white/5 p-1 transition-transform group-hover:scale-105">
          <img
            src={LOGO_SRC}
            alt="GCP Certified Logo"
            className="h-full w-full object-contain"
          />
        </div>
        <div className="space-y-1">
          <h3 className="font-semibold text-lg leading-tight group-hover:text-primary transition-colors">
            {exam.title}
          </h3>
          <div className="flex gap-2">
             {/* Proバッジ: GCP Yellowを使用 */}
            {exam.isPro ? (
              <Badge variant="outline" className="border-gcp-yellow/30 bg-gcp-yellow/10 text-gcp-yellow text-[10px] px-2 py-0 h-5">
                Professional
              </Badge>
            ) : (
              <Badge variant="outline" className="text-muted-foreground text-[10px] px-2 py-0 h-5">
                Foundational
              </Badge>
            )}
          </div>
        </div>
      </CardHeader>

      <CardContent className="pb-4">
        <div className="flex items-center gap-4 text-xs text-muted-foreground mt-2">
          <div className="flex items-center gap-1">
            <BarChart3 className="w-3 h-3" />
            <span>{exam.domainCount} Domains</span>
          </div>
          <div className="flex items-center gap-1">
             <span>•</span>
             <span>50 Questions</span>
          </div>
        </div>
      </CardContent>

      <CardFooter className="pt-0">
        <Button
          className="w-full gap-2 bg-primary hover:bg-primary/90 text-primary-foreground font-medium shadow-md shadow-primary/20"
          variant={exam.isPro ? "outline" : "default"}
        >
          {exam.isPro ? (
            <>
              <Lock className="w-4 h-4" /> 
              <span>Pro Plan Required</span>
            </>
          ) : (
            <>
              <Play className="w-4 h-4 fill-current" />
              <span>Start Exam</span>
            </>
          )}
        </Button>
      </CardFooter>
      
      {/* ホバー時のアクセントバー */}
      <div className="absolute top-0 left-0 w-1 h-full bg-primary scale-y-0 group-hover:scale-y-100 transition-transform duration-300 origin-bottom" />
    </Card>
  );
}
