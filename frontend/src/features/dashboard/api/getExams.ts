import useSWR from 'swr';
import { useClientApiClient } from '@/hooks/useClientApiClient';
import type { Exam } from '@/types/api';

export const useExams = () => {
  const client = useClientApiClient();
  
  const fetcher = (url: string) => client.get<Exam[]>(url).then(res => res.data);

  const { data, error, isLoading } = useSWR('/exams', fetcher);

  return {
    exams: data,
    isLoading,
    isError: error,
  };
};
