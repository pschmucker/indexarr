import styles from '../styles/topbar.module.css';

interface TopbarProps {
  showBack: boolean;
  breadcrumb: string;
  onBack: () => void;
}

export const Topbar = ({ showBack, breadcrumb, onBack }: TopbarProps) => {
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '12px', padding: '0 20px', height: '56px', background: 'var(--color-background-primary)', borderBottom: '0.5px solid var(--color-border-tertiary)' }}>
      {showBack && (
        <>
          <button className={styles['back-btn']} onClick={onBack}>
            <svg className={styles['back-icon']} viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5">
              <path d="M10 12L6 8l4-4" />
            </svg>
            Retour
          </button>
          <div className={styles.separator} />
        </>
      )}

      {breadcrumb && (
        <div className={styles.breadcrumb}>
          {breadcrumb}
        </div>
      )}

      <div className={styles['search-container']}>
        <svg className={styles['search-icon']} viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5">
          <circle cx="7" cy="7" r="4.5" />
          <path d="M10.5 10.5l2.5 2.5" />
        </svg>
        <input
          type="text"
          className={styles['search-input']}
          placeholder="Rechercher…"
          onFocus={(e) => {
            const shortcut = e.currentTarget.parentElement?.querySelector('[data-shortcut]');
            if (shortcut) shortcut.style.display = 'none';
          }}
          onBlur={(e) => {
            const shortcut = e.currentTarget.parentElement?.querySelector('[data-shortcut]');
            if (shortcut) shortcut.style.display = 'block';
          }}
        />
        <span className={styles.shortcut} data-shortcut>
          /
        </span>
      </div>
    </div>
  );
};
