import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import ACE from '@/assets/certifications/ACE.png';
import ADP from '@/assets/certifications/ADP.png';
import AWA from '@/assets/certifications/AWA.png';
import PCA from '@/assets/certifications/PCA.png';
import PCD from '@/assets/certifications/PCD.png';
import PCDE from '@/assets/certifications/PCDE.png';
import PCDOE from '@/assets/certifications/PCDOE.png';
import PCNE from '@/assets/certifications/PCNE.png';
import PCSE from '@/assets/certifications/PCSE.png';
import PDE from '@/assets/certifications/PDE.png';
import PMLE from '@/assets/certifications/PMLE.png';
import PSOE from '@/assets/certifications/PSOE.png';

const certifications = [
  { src: ACE, alt: 'Associate Cloud Engineer', id: 'ace' },
  { src: PCA, alt: 'Professional Cloud Architect', id: 'pca' },
  { src: PCD, alt: 'Professional Cloud Developer', id: 'pcd' },
  { src: PDE, alt: 'Professional Data Engineer', id: 'pde' },
  { src: PCDE, alt: 'Professional Cloud Database Engineer', id: 'pcde' },
  { src: PCNE, alt: 'Professional Cloud Network Engineer', id: 'pcne' },
  { src: PCSE, alt: 'Professional Cloud Security Engineer', id: 'pcse' },
  { src: PMLE, alt: 'Professional Machine Learning Engineer', id: 'pmle' },
  { src: PCDOE, alt: 'Professional Cloud DevOps Engineer', id: 'pcdoe' },
  { src: AWA, alt: 'Professional Google Workspace Administrator', id: 'awa' },
  { src: ADP, alt: 'Professional Cloud Digital Leader', id: 'adp' },
  { src: PSOE, alt: 'Professional Cloud Security Operations', id: 'psoe' },
];

export const CertificationsSection = () => {
  return (
    <section className="py-20 relative z-10">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-2xl md:text-3xl font-bold text-white mb-4">
            対応している認定資格
          </h2>
          <p className="text-gray-400">
            学習したい試験を選択してください
          </p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {certifications.map((cert) => (
            <Card key={cert.id} className="bg-card/5 border-white/10 backdrop-blur-sm hover:bg-card/10 transition-all duration-300 group">
              <CardHeader className="p-6 pb-2 flex flex-row items-center gap-4">
                <div className="w-16 h-16 shrink-0 flex items-center justify-center bg-white/5 rounded-lg p-2">
                  <img
                    src={cert.src}
                    alt={cert.alt}
                    className="w-full h-full object-contain"
                  />
                </div>
                <h3 className="font-semibold text-white leading-tight group-hover:text-blue-400 transition-colors">
                  {cert.alt}
                </h3>
              </CardHeader>
              <CardContent className="px-6 py-2">
                <div className="text-sm text-gray-400">
                  <span className="inline-block w-2 h-2 rounded-full bg-green-500 mr-2"></span>
                  利用可能
                </div>
              </CardContent>
              <CardFooter className="p-6 pt-2">
                <Button 
                  className="w-full btn-gradient text-white border-0"
                  onClick={() => window.location.href = `/exams/${cert.id}`}
                >
                  問題集を解く
                </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
};
