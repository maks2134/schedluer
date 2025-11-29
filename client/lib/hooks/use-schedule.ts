'use client';

import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api';
import type { ScheduleResponse } from '@/types/api';

export function useGroupSchedule(groupNumber: string | null, useCache = true) {
  const [data, setData] = useState<ScheduleResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!groupNumber) {
      return;
    }

    setLoading(true);
    setError(null);
    setData(null);

    apiClient
      .getGroupSchedule(groupNumber, useCache)
      .then(setData)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [groupNumber, useCache]);

  const refresh = async () => {
    if (!groupNumber) return;
    setLoading(true);
    setError(null);
    try {
      await apiClient.refreshGroupSchedule(groupNumber);
      const newData = await apiClient.getGroupSchedule(groupNumber, false);
      setData(newData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  return { data, loading, error, refresh };
}

export function useEmployeeSchedule(urlId: string | null, useCache = true) {
  const [data, setData] = useState<ScheduleResponse | null>(null);
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
      .getEmployeeSchedule(urlId, useCache)
      .then(setData)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [urlId, useCache]);

  const refresh = async () => {
    if (!urlId) return;
    setLoading(true);
    setError(null);
    try {
      await apiClient.refreshEmployeeSchedule(urlId);
      const newData = await apiClient.getEmployeeSchedule(urlId, false);
      setData(newData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  return { data, loading, error, refresh };
}

