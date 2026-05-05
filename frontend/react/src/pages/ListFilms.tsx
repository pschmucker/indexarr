import { useEffect, useState } from 'react';
import { Movie, PaginatedResponse } from '../types';
import { apiClient } from '../api/client';
import { MovieCard } from '../components/MovieCard';
import { StatCard } from '../components/StatCard';
import comStyles from '../styles/components.module.css';

interface ListFilmsProps {
  onSelectMovie: (id: number) => void;
}

export const ListFilms = ({ onSelectMovie }: ListFilmsProps) => {
  const [movies, setMovies] = useState<Movie[]>([]);
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState({ available: 0, total: 0, diskSpace: 0, fourK: 0 });

  useEffect(() => {
    const fetchMovies = async () => {
      setLoading(true);
      try {
        const response = await apiClient.getMovies(1, 50);
        setMovies(response.data);
        // Calculate stats
        const available = response.data.filter((m) => m.status === 'available').length;
        const diskSpace = response.data.reduce((sum, m) => sum + (m.fileSize || 0), 0) / (1024 * 1024 * 1024);
        const fourK = response.data.filter((m) => m.mediaInfo?.videoTracks[0]?.resolution.includes('3840')).length;
        setStats({ available, total: response.data.length, diskSpace, fourK });
      } catch (error) {
        console.error('Failed to fetch movies:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchMovies();
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
        <StatCard label="Films" value={stats.total} subLabel={`${stats.available} disponibles`} />
        <StatCard label="Espace" value={`${stats.diskSpace.toFixed(1)} Go`} subLabel="moy. disque" />
        <StatCard label="4K UHD" value={stats.fourK} subLabel={`${Math.round((stats.fourK / stats.total) * 100)}%`} />
        <StatCard label="Problèmes" value="0" subLabel="fichiers manquants" />
      </div>

      {/* Grid */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(148px, 1fr))', gap: '12px' }}>
        {movies.map((movie) => (
          <MovieCard key={movie.id} movie={movie} onClick={() => onSelectMovie(movie.id)} />
        ))}
      </div>
    </div>
  );
};
