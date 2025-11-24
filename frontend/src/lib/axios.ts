import Axios from 'axios';

// Base API URL from environment variables
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const axios = Axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for adding auth token
axios.interceptors.request.use(
  (config) => {
    // TODO: Add Firebase Auth token
    // const token = getAuthToken();
    // if (token) {
    //   config.headers.Authorization = `Bearer ${token}`;
    // }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Response interceptor for error handling
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle common errors
    if (error.response?.status === 401) {
      // Handle unauthorized access
      console.error('Unauthorized access - redirecting to login');
    }
    return Promise.reject(error);
  },
);
