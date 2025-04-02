import { Link, useLocation } from 'react-router';
import { Settings } from 'lucide-react';
import { Button } from '@/components/ui/button.tsx';
import { useMemo } from 'react';

export default function Header() {
  const pageTitles = useMemo(() => ({
    '/': 'Home',
    '/settings': 'Settings',
  }), []);
  const location = useLocation();

  return (
    <header className="grid grid-cols-3 items-center p-4 border-b-2">
      <div></div>
      <div>
        <h2 className="text-xl font-bold text-center">
          {/* @ts-expect-error Types */}
          {pageTitles[location.pathname] || ''}
        </h2>
      </div>
      <div className="flex justify-end items-center gap-4">
        <Button variant="outline" asChild>
          <Link to="/settings">
            <Settings /> Settings
          </Link>
        </Button>
      </div>
    </header>
  );
}