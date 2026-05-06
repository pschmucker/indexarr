import styles from '../styles/sidebar.module.css';
import { Page } from '../hooks/useAppContext';

interface SidebarProps {
  activeNav: string;
  onNavClick: (page: Page, id?: number) => void;
}

export const Sidebar = ({ activeNav, onNavClick }: SidebarProps) => {
  return (
    <div className={styles.sidebar}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
        <div className={styles['logo-mark']}>
          <svg viewBox="0 0 14 14" style={{ width: '13px', height: '13px', fill: 'white' }}>
            <path d="M2 11L7 3L12 11Z" />
          </svg>
        </div>
        <span className={styles['logo-name']}>Indexarr</span>
      </div>

      <nav className={styles.nav}>
        <div className={styles['nav-group']}>Librairie</div>

        <div
          className={`${styles['nav-item']} ${activeNav === 'list-films' ? styles.active : ''}`}
          onClick={() => onNavClick('list-films')}
        >
          <svg className={styles['nav-icon']} viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5">
            <rect x="2" y="3" width="12" height="10" rx="1.5" />
            <path d="M5 3v10M11 3v10M2 7h12" />
          </svg>
          Films
          <span className={styles['nav-badge']}>20</span>
        </div>

        <div
          className={`${styles['nav-item']} ${activeNav === 'list-series' ? styles.active : ''}`}
          onClick={() => onNavClick('list-series')}
        >
          <svg className={styles['nav-icon']} viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5">
            <rect x="2" y="2" width="12" height="12" rx="1.5" />
            <path d="M2 6h12M6 6v8" />
          </svg>
          Séries
          <span className={styles['nav-badge']}>10</span>
        </div>

        <div className={styles['nav-item']}>
          <svg className={styles['nav-icon']} viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5">
            <circle cx="8" cy="8" r="5" />
            <path d="M8 5v3l2 2" />
          </svg>
          Récents
        </div>

        <div className={styles['nav-group']} style={{ marginTop: '6px' }}>
          Analyse
        </div>

        <div className={styles['nav-item']}>
          <svg className={styles['nav-icon']} viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5">
            <path d="M2 12l4-5 3 3 3-4 2 2" />
          </svg>
          Statistiques
        </div>

        <div className={styles['nav-item']}>
          <svg className={styles['nav-icon']} viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5">
            <circle cx="8" cy="8" r="6" />
            <path d="M8 5v4M8 11h.01" />
          </svg>
          Problèmes
          <span className={styles['nav-badge']} style={{ background: '#FCEBEB', color: '#791F1F', borderColor: '#F09595' }}>
            4
          </span>
        </div>
      </nav>

      <div className={styles.footer}>
        <div className={styles['status-dot']} />
        <span>Système opérationnel</span>
      </div>
    </div>
  );
};
