import { useNavigate } from 'react-router-dom';
import { loginWithGoogle } from '@/features/auth/api/login';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Github, Apple, Chrome } from 'lucide-react';
import logo from '@/assets/logo.png';

export const LoginPage = () => {
  const navigate = useNavigate();

  const handleGoogleLogin = async () => {
    try {
      await loginWithGoogle();
      navigate('/');
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  const handleGithubLogin = () => {
    // TODO: Implement Github login
    console.log('Github login not implemented yet');
  };

  const handleAppleLogin = () => {
    // TODO: Implement Apple login
    console.log('Apple login not implemented yet');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-background relative overflow-hidden">
      {/* Background Effects */}
      <div className="absolute top-0 left-0 w-full h-full overflow-hidden -z-10 pointer-events-none">
        <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-[100px]" />
        <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-blue-600/20 rounded-full blur-[100px]" />
      </div>

      <Card className="w-full max-w-md mx-4 border-white/10 bg-black/40 backdrop-blur-xl shadow-2xl">
        <CardHeader className="text-center space-y-2">
          <div className="mx-auto w-12 h-12 mb-4">
            <img src={logo} alt="Nearline Logo" className="w-full h-full object-contain" />
          </div>
          <CardTitle className="text-2xl font-bold text-white tracking-tight">
            Welcome back
          </CardTitle>
          <CardDescription className="text-gray-400">
            Sign in to your account to continue
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <Button 
            variant="outline" 
            className="w-full h-12 bg-white/5 border-white/10 hover:bg-white/10 hover:text-white text-gray-200 justify-start px-4 gap-3 transition-all duration-300"
            onClick={handleGoogleLogin}
          >
            <Chrome className="w-5 h-5" />
            <span>Continue with Google</span>
          </Button>

          <Button 
            variant="outline" 
            className="w-full h-12 bg-white/5 border-white/10 hover:bg-white/10 hover:text-white text-gray-200 justify-start px-4 gap-3 transition-all duration-300 opacity-50 cursor-not-allowed"
            onClick={handleGithubLogin}
            disabled
          >
            <Github className="w-5 h-5" />
            <span>Continue with GitHub (Coming Soon)</span>
          </Button>

          <Button 
            variant="outline" 
            className="w-full h-12 bg-white/5 border-white/10 hover:bg-white/10 hover:text-white text-gray-200 justify-start px-4 gap-3 transition-all duration-300 opacity-50 cursor-not-allowed"
            onClick={handleAppleLogin}
            disabled
          >
            <Apple className="w-5 h-5" />
            <span>Continue with Apple (Coming Soon)</span>
          </Button>

          <div className="relative my-6">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t border-white/10" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-black/40 px-2 text-gray-500">
                Protected by Nearline Security
              </span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};
