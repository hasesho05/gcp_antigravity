import { axios } from '@/lib/axios';
import type { AxiosInstance } from 'axios';

export const useClientApiClient = (): AxiosInstance => {
  // 将来的にここでトークンの取得やヘッダーの付与を行う
  // const { token } = useAuth();
  // if (token) {
  //   axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
  // }
  return axios;
};
