import { BrainCircuit, BookOpen, Smartphone } from 'lucide-react';

const features = [
  {
    title: "苦手を克服する、画期的な学習システム",
    description: "学習履歴から正答率の低い分野を自動で抽出。AIがあなたの弱点を分析し、最適な復習タイミングを提案。記憶の定着を徹底的にサポートします。",
    icon: BrainCircuit,
    color: "text-blue-400",
    bg: "bg-blue-500/10",
    border: "border-blue-500/20"
  },
  {
    title: "本質を突く詳細な解説",
    description: "現役GCPエンジニアが執筆した、図解豊富な解説。単なる正誤判定だけでなく、「なぜそうなるのか」という技術の裏側まで深く理解できます。",
    icon: BookOpen,
    color: "text-green-400",
    bg: "bg-green-500/10",
    border: "border-green-500/20"
  },
  {
    title: "いつでもどこでも、スマートに学習",
    description: "PC、タブレット、スマートフォンに完全対応。通勤中の電車内や、ちょっとした待ち時間も、貴重なスキルアップの時間に変えます。",
    icon: Smartphone,
    color: "text-yellow-400",
    bg: "bg-yellow-500/10",
    border: "border-yellow-500/20"
  },
];

export const FeaturesSection = () => {
  return (
    <section className="py-24 relative overflow-hidden">
        {/* Background Gradients */}
        <div className="absolute top-1/2 left-0 w-[500px] h-[500px] bg-blue-600/10 rounded-full blur-[100px] -translate-y-1/2 -translate-x-1/2 pointer-events-none" />
        <div className="absolute bottom-0 right-0 w-[500px] h-[500px] bg-purple-600/10 rounded-full blur-[100px] translate-y-1/3 translate-x-1/3 pointer-events-none" />

      <div className="container mx-auto px-4 relative z-10">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-6">
            Nearlineが選ばれ続ける理由
          </h2>
          <p className="text-gray-400 max-w-2xl mx-auto text-lg">
            合格に必要なのは、ただ問題を解くことではありません。<br className="hidden md:block" />
            効率的かつ本質的な学習体験を提供します。
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-8">
          {features.map((feature, index) => (
            <div 
              key={index} 
              className={`group relative bg-card/5 border ${feature.border} rounded-2xl p-8 backdrop-blur-sm hover:bg-card/10 transition-all duration-300 hover:-translate-y-1`}
            >
              <div className={`w-16 h-16 ${feature.bg} rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300`}>
                <feature.icon className={`w-8 h-8 ${feature.color}`} />
              </div>
              <h3 className="text-xl font-bold text-white mb-4">
                {feature.title}
              </h3>
              <p className="text-gray-400 leading-relaxed">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};
