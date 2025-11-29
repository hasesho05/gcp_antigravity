import { axios } from '@/lib/axios';
import { auth } from '@/lib/firebase';
import type { AxiosInstance } from 'axios';
import { useEffect } from 'react';

export const useClientApiClient = (): AxiosInstance => {
  useEffect(() => {
    const interceptorId = axios.interceptors.request.use(
      async (config) => {
        const user = auth.currentUser;
        if (user) {
          const token = await user.getIdToken();
          if (token) {
            config.headers.Authorization = `Bearer ${token}`;
          }
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    return () => {
      axios.interceptors.request.eject(interceptorId);
    };
  }, []);

  return axios;
};
