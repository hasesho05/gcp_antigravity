import { Header } from '@/components/Header';
import { HeroSection } from '@/components/HeroSection';
import { FeaturesSection } from '@/components/FeaturesSection';
import { CertificationsSection } from '@/components/CertificationsSection';

export const LandingPage = () => {
  return (
    <div className="min-h-screen text-foreground overflow-x-hidden relative">
      <Header />
      <main>
        <HeroSection />
        <FeaturesSection />
        <CertificationsSection />
      </main>
    </div>
  );
};
