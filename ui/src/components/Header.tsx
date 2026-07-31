import { Link, useLocation } from 'react-router-dom';

export function Header() {
  const location = useLocation();

  return (
    <header className="w-full flex items-center justify-between py-4 px-8 border-b border-border bg-background">
      <Link to="/" className="text-xl font-bold text-primary">
        smriti
      </Link>
      <nav className="flex gap-6">
        <Link
          to="/"
          className={`${
            location.pathname === '/'
              ? 'text-primary underline underline-offset-4'
              : 'text-foreground hover:text-primary transition-colors'
          }`}
        >
          Home
        </Link>
        <Link
          to="/about"
          className={`${
            location.pathname === '/about'
              ? 'text-primary underline underline-offset-4'
              : 'text-foreground hover:text-primary transition-colors'
          }`}
        >
          About
        </Link>
      </nav>
    </header>
  );
}
