import { Movie } from '../types';
import comStyles from '../styles/components.module.css';

interface MovieCardListProps {
  movie: Movie;
  onClick: () => void;
}

export const MovieCardList = ({ movie, onClick }: MovieCardListProps) => {
  const initials = movie.title
    .split(' ')
    .map((word) => word[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);

  const statusColor = movie.status === 'available' ? '#1D9E75' : movie.status === 'missing' ? '#E24B4A' : '#EF9F27';

  return (
    <div className={comStyles['card-list']} onClick={onClick}>
      {movie.poster ? (
        <div className={comStyles['card-list-poster']}>
          <img
            src={movie.poster}
            alt={movie.title}
            width="100%"
            height="100%"
            style={{ objectFit: 'contain' }}
            className={comStyles['card-list-poster-img']}
          />
          <div className={comStyles['card-list-poster-status']} style={{ background: statusColor }} />
        </div>
      ) : (
        <div className={comStyles['card-list-poster']}>
          <div className={comStyles['card-list-poster-initial']}>
            {initials}
          </div>
          <div className={comStyles['card-list-poster-title']}>
            {movie.title}
          </div>
          <div className={comStyles['card-list-poster-status']} style={{ background: statusColor }} />
        </div>
      )}

      <div className={comStyles['card-list-content']}>
        <div className={comStyles['card-list-title']}>
          {movie.title}
        </div>
        <div className={comStyles['card-list-meta']}>
          <span>{movie.year}</span>
          <span>·</span>
          <span>{movie.duration} min</span>
          <span>·</span>
          <span>{movie.genres}</span>
          <span>·</span>
          <span style={{ color: statusColor, fontWeight: 500 }}>
            {movie.status === 'available' ? 'Disponible' : movie.status === 'missing' ? 'Manquant' : 'Problème'}
          </span>
        </div>
        <div className={comStyles['card-list-badges']}>
          {movie.mediaInfo?.videoTracks?.[0]?.resolution.includes('x2160') && <span className={comStyles['badge-4k']}>4K</span>}
          {movie.mediaInfo?.videoTracks?.[0]?.hdr.includes('Dolby') && <span className={comStyles['badge-dv']}>DV</span>}
          {movie.mediaInfo?.videoTracks?.[0]?.hdr.includes('HDR10+') && <span className={comStyles['badge-hdr']}>HDR10+</span>}
          {movie.mediaInfo?.videoTracks?.[0]?.hdr.includes('HDR10') && !movie.mediaInfo?.videoTracks?.[0]?.hdr.includes('HDR10+') && <span className={comStyles['badge-hdr']}>HDR10</span>}
          {(movie.mediaInfo?.audioTracks ?? []).find((track) => track.codec === 'TrueHD') && <span className={comStyles['badge-truehd']}>TrueHD</span>}
          {(movie.mediaInfo?.audioTracks ?? []).find((track) => track.codec === 'E-AC-3') && <span className={comStyles['badge-ddplus']}>DD+</span>}
          {(movie.mediaInfo?.audioTracks ?? []).find((track) => track.codec.includes('Atmos')) && <span className={comStyles['badge-atmos']}>Atmos</span>}
          {(movie.mediaInfo?.audioTracks ?? []).find((track) => track.codec === 'DTS') && <span className={comStyles['badge-dts']}>DTS</span>}
          {movie.mediaInfo?.videoTracks?.[0]?.codec && <span className={comStyles['badge-codec']}>{movie.mediaInfo.videoTracks?.[0]?.codec}</span>}
          {movie.status === 'missing' && <span className={comStyles['badge-missing']}>Manquant</span>}
        </div>
      </div>
    </div>
  );
};
