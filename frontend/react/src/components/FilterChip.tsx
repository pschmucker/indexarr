import styles from '../styles/components.module.css';

interface FilterChipProps {
  label: string;
  active: boolean;
  count?: number;
  onClick: () => void;
}

export const FilterChip = ({ label, active, count, onClick }: FilterChipProps) => {
  return (
    <div
      className={`${styles['filter-chip']} ${active ? styles['filter-chip-active'] : ''}`}
      onClick={onClick}
    >
      {label}
      {count !== undefined && count > 0 && (
        <span className={styles['filter-chip-badge']}>{count}</span>
      )}
    </div>
  );
};
