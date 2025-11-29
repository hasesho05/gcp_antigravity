import useSWR from 'swr';
import { useClientApiClient } from '@/hooks/useClientApiClient';

export interface ExamSetStats {
  examSetId: string;
  totalAttempts: number;
  latestAttempt?: {
    status: 'in_progress' | 'completed' | 'paused';
    progress: number; // 0-100
    score?: number;
    lastAccessedAt: string;
  };
}

// Dummy data for development
const dummyStats: ExamSetStats[] = [
  {
    examSetId: 'set-1',
    totalAttempts: 3,
    latestAttempt: {
      status: 'in_progress',
      progress: 50,
      lastAccessedAt: new Date().toISOString(),
    },
  },
  {
    examSetId: 'set-2',
    totalAttempts: 0,
  },
  {
    examSetId: 'set-3',
    totalAttempts: 1,
    latestAttempt: {
      status: 'completed',
      progress: 100,
      score: 85,
      lastAccessedAt: new Date(Date.now() - 86400000).toISOString(),
    },
  },
];

export const useExamSetStats = (examId: string | undefined) => {
  const client = useClientApiClient();
  
  // Placeholder for future API endpoint
  const fetcher = (url: string) => client.get<ExamSetStats[]>(url).then(res => res.data);

  const { data } = useSWR(examId ? `/exams/${examId}/sets/stats` : null, fetcher);

  return {
    stats: data || dummyStats,
    isLoading: false,
    isError: null,
  };
};
