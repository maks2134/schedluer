'use client';

import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api';
import type { EmployeeListItem } from '@/types/api';

export function useEmployees(useCache = true) {
  const [data, setData] = useState<EmployeeListItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    setError(null);

    apiClient
      .getAllEmployees(useCache)
      .then(setData)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [useCache]);

  const refresh = async () => {
    setLoading(true);
    setError(null);
    try {
      await apiClient.refreshEmployees();
      const newData = await apiClient.getAllEmployees(false);
      setData(newData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  return { data, loading, error, refresh };
}

export function useEmployee(urlId: string | null) {
  const [data, setData] = useState<EmployeeListItem | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!urlId) {
      return;
    }

    setLoading(true);
    setError(null);
    setData(null);

    apiClient
      .getEmployeeByUrlId(urlId)
      .then(setData)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [urlId]);

  return { data, loading, error };
}

