import useSWR from 'swr';
import { useClientApiClient } from '@/hooks/useClientApiClient';
import type { Exam } from '@/types/api';

export const useExam = (examId: string | undefined) => {
  const client = useClientApiClient();
  
  const fetcher = (url: string) => client.get<Exam>(url).then(res => res.data);

  const { data } = useSWR(examId ? `/exams/${examId}` : null, fetcher);

  // Dummy data for development
  const dummyExam: Exam = {
    id: 'pcd',
    code: 'PCD',
    name: 'Professional Cloud Developer',
    description: 'Google Cloud Certified - Professional Cloud Developer 認定試験は、スケーラブルで可用性が高く、信頼性に優れたクラウドネイティブ アプリケーションを構築、デプロイ、管理する能力を評価します。',
    imageUrl: 'https://www.gstatic.com/images/branding/product/2x/google_cloud_64dp.png', // Placeholder
    createdAt: new Date().toISOString(),
  };

  return {
    exam: data || dummyExam, // Fallback to dummy data
    isLoading: false, // Force loading to false for immediate display
    isError: null,
  };
};
