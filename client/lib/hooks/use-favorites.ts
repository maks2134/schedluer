'use client';

import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api';
import type { FavoriteGroup } from '@/types/api';

const USER_ID = 'default';

export function useFavorites() {
  const [favorites, setFavorites] = useState<FavoriteGroup[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadFavorites();
  }, []);

  const loadFavorites = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await apiClient.getAllFavorites(USER_ID);
      setFavorites(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  const addFavorite = async (groupNumber: string) => {
    try {
      await apiClient.addFavorite(groupNumber, USER_ID);
      await loadFavorites();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
      throw err;
    }
  };

  const removeFavorite = async (groupNumber: string) => {
    try {
      await apiClient.removeFavorite(groupNumber, USER_ID);
      await loadFavorites();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
      throw err;
    }
  };

  const isFavorite = (groupNumber: string): boolean => {
    if (!favorites || !Array.isArray(favorites)) {
      return false;
    }
    return favorites.some((fav) => fav.group_number === groupNumber);
  };

  const toggleFavorite = async (groupNumber: string) => {
    const currentlyFavorite = isFavorite(groupNumber);
    
    // Оптимистичное обновление
    if (currentlyFavorite) {
      setFavorites((prev) => {
        if (!Array.isArray(prev)) return [];
        return prev.filter((fav) => fav.group_number !== groupNumber);
      });
    } else {
      const newFavorite: FavoriteGroup = {
        id: '',
        group_number: groupNumber,
        user_id: USER_ID,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      };
      setFavorites((prev) => {
        if (!Array.isArray(prev)) return [newFavorite];
        return [...prev, newFavorite];
      });
    }

    try {
      if (currentlyFavorite) {
        await removeFavorite(groupNumber);
      } else {
        await addFavorite(groupNumber);
      }
      // Перезагружаем для синхронизации с сервером
      await loadFavorites();
    } catch (err) {
      // Откатываем изменения при ошибке
      await loadFavorites();
      throw err;
    }
  };

  return {
    favorites,
    loading,
    error,
    addFavorite,
    removeFavorite,
    isFavorite,
    toggleFavorite,
    refresh: loadFavorites,
  };
}

