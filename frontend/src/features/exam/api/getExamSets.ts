import useSWR from 'swr';
import { useClientApiClient } from '@/hooks/useClientApiClient';
import type { ExamSet } from '@/types/api';

export const useExamSets = (examId: string | undefined) => {
  const client = useClientApiClient();
  
  const fetcher = (url: string) => client.get<ExamSet[]>(url).then(res => res.data);

  const { data } = useSWR(examId ? `/exams/${examId}/sets` : null, fetcher);

  // Dummy data for development
  const dummyExamSets: ExamSet[] = [
    {
      id: 'set-1',
      examId: examId || 'pcd',
      name: '模擬試験 1',
      description: '全分野を網羅した総合的な模擬試験です。',
      questionIds: Array(50).fill('q'),
      createdAt: new Date().toISOString(),
    },
    {
      id: 'set-2',
      examId: examId || 'pcd',
      name: '模擬試験 2',
      description: 'セキュリティとデプロイメントに重点を置いたセットです。',
      questionIds: Array(40).fill('q'),
      createdAt: new Date().toISOString(),
    },
    {
      id: 'set-3',
      examId: examId || 'pcd',
      name: '模擬試験 3',
      description: '応用的なトピックとケーススタディを含みます。',
      questionIds: Array(60).fill('q'),
      createdAt: new Date().toISOString(),
    },
  ];

  return {
    examSets: data || dummyExamSets, // Fallback to dummy data
    isLoading: false,
    isError: null,
  };
};
