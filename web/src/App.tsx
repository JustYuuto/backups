import Header from '@/components/Header.tsx';
import { Outlet } from 'react-router';

export default function App() {
  return (
    <>
      <Header />
      <div className="container mx-auto py-4 px-2">
        <Outlet />
      </div>
    </>
  );
}