import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { LandingPage } from '@/routes/LandingPage';
import { LoginPage } from '@/features/auth/routes/LoginPage';
import { ExamDetailPage } from '@/features/exam/routes/ExamDetailPage';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/exams/:examId" element={<ExamDetailPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
