import { Link } from 'react-router-dom';
import logo from '@/assets/logo.png';
import { Button } from '@/components/ui/button';

export const Header = () => {
  return (
    <header className="fixed top-0 w-full z-50 bg-background/80 backdrop-blur-md border-b border-white/10">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2">
          <img src={logo} alt="Nearline Logo" className="h-8 w-auto" />
          <span className="text-xl font-bold tracking-tight text-white">NEARLINE</span>
        </Link>
        
        <nav className="hidden md:flex items-center gap-8">
          <Link to="/exams" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">
            Exams
          </Link>
          <Link to="/about" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">
            About
          </Link>
          <Link to="/pricing" className="text-sm font-medium text-gray-300 hover:text-white transition-colors">
            Pricing
          </Link>
        </nav>

        <div className="flex items-center gap-4">
          <Link to="/login" className="text-sm font-medium text-white hover:text-primary transition-colors">
            Log in
          </Link>
          <Button className="btn-gradient text-white border-0 font-semibold">
            Sign up
          </Button>
        </div>
      </div>
    </header>
  );
};
