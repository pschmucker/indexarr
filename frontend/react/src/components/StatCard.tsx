import comStyles from '../styles/components.module.css';

interface StatCardProps {
  label: string;
  value: string | number;
  subLabel?: string;
  error?: boolean;
}

export const StatCard = ({ label, value, subLabel, error }: StatCardProps) => {
  return (
    <div className={comStyles.stat}>
      <div className={comStyles['stat-label']}>{label}</div>
      <div className={comStyles['stat-value']} style={{ color: error ? '#E24B4A' : 'var(--color-text-primary)' }}>
        {value}
      </div>
      {subLabel && <div className={comStyles['stat-sub']}>{subLabel}</div>}
    </div>
  );
};
