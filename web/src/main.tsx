import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import App from './App.tsx';
import { BrowserRouter, Route, Routes } from 'react-router';
import Settings from '@/routes/settings.tsx';
import Home from '@/routes/home.tsx';
import { Toaster } from '@/components/ui/sonner.tsx';

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Toaster />
    <BrowserRouter>
      <Routes>
        <Route element={<App />}>
          <Route path="/" element={<Home />} />
          <Route path="/settings" element={<Settings />} />
        </Route>
      </Routes>
    </BrowserRouter>
  </StrictMode>,
);
