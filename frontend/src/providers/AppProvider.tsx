import type React from 'react';
import { SWRConfig } from 'swr';
import { swrConfig } from '@/lib/swr';

interface AppProviderProps {
  children: React.ReactNode;
}

export const AppProvider: React.FC<AppProviderProps> = ({ children }) => {
  return <SWRConfig value={swrConfig}>{children}</SWRConfig>;
};
