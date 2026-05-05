import { useContext, ReactNode } from 'react';
import { AppContext } from '../hooks/useAppContext';

interface LayoutProps {
  children: ReactNode;
  currentPage: string;
}

export const Layout = ({ children }: LayoutProps) => {
  const context = useContext(AppContext);
  if (!context) return null;

  const { currentPage, goToPage, goBack, history } = context;

  return (
    <div style={{ display: 'flex', height: '100vh', width: '100%', background: 'var(--color-background-tertiary)' }}>
      {/* Sidebar will be added here */}
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        {/* Topbar will be added here */}
        <div style={{ flex: 1, overflow: 'y', background: 'var(--color-background-primary)' }}>{children}</div>
      </div>
    </div>
  );
};
