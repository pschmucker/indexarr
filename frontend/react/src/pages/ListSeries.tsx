import { useEffect, useState } from 'react';
import { Series, PaginatedResponse } from '../types';
import { apiClient } from '../api/client';
import { SeriesCard } from '../components/SeriesCard';
import { StatCard } from '../components/StatCard';

interface ListSeriesProps {
  onSelectSeries: (id: number) => void;
}

export const ListSeries = ({ onSelectSeries }: ListSeriesProps) => {
  const [series, setSeries] = useState<Series[]>([]);
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState({ complete: 0, total: 0, episodes: 0, diskSpace: 0 });

  useEffect(() => {
    const fetchSeries = async () => {
      setLoading(true);
      try {
        const response = await apiClient.getSeries(1, 50);
        setSeries(response.data);
        // Calculate stats
        const complete = response.data.filter((s) => s.status === 'complete').length;
        const episodes = response.data.reduce((sum, s) => sum + s.episodeCount, 0);
        const diskSpace = response.data.reduce((sum, s) => sum + (s.fileSize || 0), 0) / (1024 * 1024 * 1024 * 1024);
        setStats({ complete, total: response.data.length, episodes, diskSpace });
      } catch (error) {
        console.error('Failed to fetch series:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchSeries();
  }, []);

  return (
    <div style={{ padding: '16px 20px' }}>
      {/* Filters */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '6px', marginBottom: '16px', paddingBottom: '8px', borderBottom: '0.5px solid var(--color-border-tertiary)' }}>
        <span style={{ fontSize: '11px', color: 'var(--color-text-tertiary)', marginRight: '2px' }}>Filtres</span>
        <div style={{ border: '0.5px solid var(--color-border-tertiary)', background: 'var(--color-background-secondary)', borderRadius: '99px', padding: '4px 10px', fontSize: '11px', color: 'var(--color-text-secondary)', cursor: 'pointer' }}>
          Statut
        </div>
      </div>

      {/* Stats */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, 1fr)', gap: '10px', marginBottom: '16px' }}>
        <StatCard label="Séries" value={stats.total} subLabel={`${stats.complete} complètes`} />
        <StatCard label="Épisodes" value={stats.episodes} subLabel="total" />
        <StatCard label="Espace" value={`${stats.diskSpace.toFixed(1)} To`} subLabel="moy. par ep." />
        <StatCard label="Problèmes" value="0" subLabel="épisodes manquants" />
      </div>

      {/* Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(148px, 1fr))', gap: '12px' }}>
        {series.map((s) => (
          <SeriesCard key={s.id} series={s} onClick={() => onSelectSeries(s.id)} />
        ))}
      </div>
    </div>
  );
};
