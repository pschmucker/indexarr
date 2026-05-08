import { createContext, useState, ReactNode } from 'react';

export type Page = 'list-films' | 'list-series' | 'detail-movie' | 'detail-series';

interface AppContextType {
  currentPage: Page;
  selectedId: number | null;
  goToPage: (page: Page, id?: number) => void;
  goBack: () => void;
  history: Page[];
  isDark: boolean;
  toggleTheme: () => void;
}

export const AppContext = createContext<AppContextType | undefined>(undefined);

interface AppContextProviderProps {
  children: ReactNode;
}

export const AppContextProvider = ({ children }: AppContextProviderProps) => {
  const [currentPage, setCurrentPage] = useState<Page>('list-films');
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [history, setHistory] = useState<Page[]>(['list-films']);
  const [isDark, setIsDark] = useState(false);

  const goToPage = (page: Page, id?: number) => {
    setCurrentPage(page);
    if (id) setSelectedId(id);
    setHistory([...history, page]);
  };

  const goBack = () => {
    if (history.length > 1) {
      const newHistory = history.slice(0, -1);
      setHistory(newHistory);
      setCurrentPage(newHistory[newHistory.length - 1]);
    }
  };

  const toggleTheme = () => {
    setIsDark(!isDark);
    document.documentElement.style.colorScheme = isDark ? 'light' : 'dark';
  };

  return (
    <AppContext.Provider value={{ currentPage, selectedId, goToPage, goBack, history, isDark, toggleTheme }}>
      {children}
    </AppContext.Provider>
  );
};
