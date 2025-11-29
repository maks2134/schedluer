'use client';

import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api';
import type { StudentGroupListItem } from '@/types/api';

export function useGroups(useCache = true) {
  const [data, setData] = useState<StudentGroupListItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    setError(null);

    apiClient
      .getAllGroups(useCache)
      .then(setData)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [useCache]);

  const refresh = async () => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.refreshGroups();
      const newData = await apiClient.getAllGroups(false);
      setData(newData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  return { data, loading, error, refresh };
}

export function useGroup(groupNumber: string | null) {
  const [data, setData] = useState<StudentGroupListItem | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!groupNumber) {
      setData(null);
      return;
    }

    setLoading(true);
    setError(null);

    apiClient
      .getGroupByNumber(groupNumber)
      .then(setData)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [groupNumber]);

  return { data, loading, error };
}

