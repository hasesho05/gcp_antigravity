import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import App from '@/App';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
  },
  // TODO: Add feature routes
  // {
  //   path: '/auth',
  //   element: <AuthLayout />,
  //   children: [
  //     { path: 'login', element: <LoginPage /> },
  //     { path: 'register', element: <RegisterPage /> },
  //   ],
  // },
  // {
  //   path: '/exam',
  //   element: <MainLayout />,
  //   children: [
  //     { path: ':examId', element: <ExamPage /> },
  //   ],
  // },
]);

export const AppRoutes = () => {
  return <RouterProvider router={router} />;
};
