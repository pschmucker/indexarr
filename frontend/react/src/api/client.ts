import { Movie, Series, Episode, PaginatedResponse, StatsResponse } from '../types/index';

const API_BASE = '/api';

export const apiClient = {
  getMovies: async (page: number = 1, pageSize: number = 50, filters: Record<string, string> = {}) => {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
      ...filters,
    });
    const response = await fetch(`${API_BASE}/movies?${params}`);
    return response.json() as Promise<PaginatedResponse<Movie>>;
  },

  getMovie: async (id: number) => {
    const response = await fetch(`${API_BASE}/movies/${id}`);
    return response.json() as Promise<Movie>;
  },

  getSeries: async (page: number = 1, pageSize: number = 50, filters: Record<string, string> = {}) => {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
      ...filters,
    });
    const response = await fetch(`${API_BASE}/series?${params}`);
    return response.json() as Promise<PaginatedResponse<Series>>;
  },

  getSeriesById: async (id: number) => {
    const response = await fetch(`${API_BASE}/series/${id}`);
    return response.json() as Promise<Series>;
  },

  getStats: async () => {
    const response = await fetch(`${API_BASE}/stats`);
    return response.json() as Promise<StatsResponse>;
  },
};
