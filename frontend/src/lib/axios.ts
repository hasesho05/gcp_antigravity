import Axios from 'axios';
import { env } from '@/config/env';

export const axios = Axios.create({
  baseURL: env.API_URL,
});

axios.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    return Promise.reject(error);
  }
);
