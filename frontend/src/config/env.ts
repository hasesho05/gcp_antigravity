// Environment variables configuration
export const env = {
  API_URL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  NODE_ENV: import.meta.env.MODE,
  isDevelopment: import.meta.env.DEV,
  isProduction: import.meta.env.PROD,
} as const;
