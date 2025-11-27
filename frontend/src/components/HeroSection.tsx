import { Button } from '@/components/ui/button';
import { ArrowRight, CheckCircle2 } from 'lucide-react';

export const HeroSection = () => {
  return (
    <section className="relative pt-32 pb-20 lg:pt-48 lg:pb-32 overflow-hidden">
      {/* Background Elements */}
      <div className="absolute top-0 left-0 w-full h-full overflow-hidden -z-10 pointer-events-none">
        <div className="absolute top-[-10%] right-[-5%] w-[500px] h-[500px] bg-blue-600/20 rounded-full blur-[100px]" />
        <div className="absolute bottom-[-10%] left-[-10%] w-[600px] h-[600px] bg-purple-600/10 rounded-full blur-[120px]" />
      </div>

      <div className="container mx-auto px-4">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left Column: Content */}
          <div className="space-y-8 animate-fade-in-up">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/5 border border-white/10 text-sm text-blue-300 backdrop-blur-sm">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-blue-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-blue-500"></span>
              </span>
              New: Professional Cloud Architect Updated
            </div>

            <h1 className="text-4xl lg:text-6xl font-bold leading-tight text-white tracking-tight">
              Google Cloudスキルを、<br />
              <span className="text-gradient">効率的に。</span><br />
              未来を拓く学習プラットフォーム。
            </h1>
            
            <p className="text-lg text-gray-400 max-w-xl leading-relaxed">
              実践的なラボと最新の教材で、あなたのキャリアを加速させます。
              Nearlineで、確かな技術力を手に入れましょう。
            </p>

            <div className="flex flex-col sm:flex-row gap-4 pt-4">
              <Button size="lg" className="btn-gradient text-white font-bold text-lg px-8 h-14 rounded-full shadow-lg shadow-blue-500/20 hover:scale-105 transition-transform">
                無料で始める
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
              <Button variant="ghost" size="lg" className="text-white hover:bg-white/10 font-medium text-lg h-14 rounded-full border border-white/10 hover:border-white/20">
                プランを見る
              </Button>
            </div>

            <div className="pt-8 flex flex-wrap items-center gap-6 text-sm text-gray-400">
              <div className="flex items-center gap-2">
                <CheckCircle2 className="w-5 h-5 text-green-500" />
                <span>Latest GCP Content</span>
              </div>
              <div className="flex items-center gap-2">
                <CheckCircle2 className="w-5 h-5 text-blue-500" />
                <span>Interactive Labs</span>
              </div>
              <div className="flex items-center gap-2">
                <CheckCircle2 className="w-5 h-5 text-yellow-500" />
                <span>Certified Instructors</span>
              </div>
            </div>
          </div>

          {/* Right Column: Visual */}
          <div className="relative lg:h-[600px] flex items-center justify-center [perspective:1000px]">
             <div className="relative w-full h-full max-w-[600px] max-h-[600px] transform hover:scale-[1.02] transition-transform duration-500">
                {/* Glow Effect */}
                <div className="absolute inset-0 bg-gradient-to-tr from-blue-500/20 via-purple-500/10 to-yellow-500/20 rounded-full blur-3xl opacity-40 animate-pulse" />
                
                <img 
                  src="/hero_gemini.png" 
                  alt="Cloud Technology Visualization" 
                  className="relative z-10 w-full h-full object-contain drop-shadow-2xl"
                />
             </div>
          </div>
        </div>
      </div>
    </section>
  );
};
