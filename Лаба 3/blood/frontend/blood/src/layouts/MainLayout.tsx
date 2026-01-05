import { ReactNode } from 'react';
import { useLocation } from 'react-router-dom';
import Header from '../components/Header';
import './MainLayout.css';

interface MainLayoutProps {
  children: ReactNode;
}

const MainLayout = ({ children }: MainLayoutProps) => {
  const location = useLocation();

  return (
    <div className="app-container">
      <Header currentPath={location.pathname} />
      <main className="main-wrapper">
        {children}
      </main>
    </div>
  );
};

export default MainLayout;
