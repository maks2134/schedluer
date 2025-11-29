'use client';

import { useState, useMemo } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Search, RefreshCw, Star } from 'lucide-react';
import type { StudentGroupListItem } from '@/types/api';
import { useFavorites } from '@/lib/hooks/use-favorites';

interface GroupSelectorProps {
  groups: StudentGroupListItem[];
  loading: boolean;
  selectedGroup: string | null;
  onSelectGroup: (groupNumber: string) => void;
  onRefresh?: () => void;
  refreshing?: boolean;
}

export function GroupSelector({
  groups,
  loading,
  selectedGroup,
  onSelectGroup,
  onRefresh,
  refreshing,
}: GroupSelectorProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const { isFavorite, toggleFavorite, loading: favoritesLoading, favorites } = useFavorites();

  const filteredGroups = useMemo(() => {
    return groups.filter((group) =>
      group.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      group.specialityName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      group.facultyName.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [groups, searchQuery]);

  // Сортируем: избранные группы сверху
  const sortedGroups = useMemo(() => {
    const favoriteGroupNumbers = new Set(
      Array.isArray(favorites) ? favorites.map((f) => f.group_number) : []
    );
    
    return [...filteredGroups].sort((a, b) => {
      const aIsFavorite = favoriteGroupNumbers.has(a.name);
      const bIsFavorite = favoriteGroupNumbers.has(b.name);
      if (aIsFavorite && !bIsFavorite) return -1;
      if (!aIsFavorite && bIsFavorite) return 1;
      return 0;
    });
  }, [filteredGroups, favorites]);

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-32" />
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {[...Array(5)].map((_, i) => (
              <Skeleton key={i} className="h-16 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle>Группы</CardTitle>
            <CardDescription>
              Найдено: {sortedGroups.length} из {groups.length}
            </CardDescription>
          </div>
          {onRefresh && (
            <Button
              variant="outline"
              size="sm"
              onClick={onRefresh}
              disabled={refreshing}
            >
              <RefreshCw className={`h-4 w-4 ${refreshing ? 'animate-spin' : ''}`} />
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Поиск по номеру группы..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9 text-sm sm:text-base w-full"
              autoComplete="off"
            />
          </div>
          <div className="max-h-[400px] sm:max-h-[600px] overflow-y-auto space-y-2 overscroll-contain">
            {sortedGroups.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                Группы не найдены
              </div>
            ) : (
              sortedGroups.map((group) => {
                const favorite = isFavorite(group.name);
                return (
                  <div
                    key={group.id}
                    className={`w-full rounded-lg border transition-all ${
                      selectedGroup === group.name
                        ? 'border-primary bg-accent'
                        : 'border-border'
                    }`}
                  >
                    <div className="flex items-start gap-2 p-3">
                      <button
                        onClick={() => onSelectGroup(group.name)}
                        className="flex-1 text-left min-w-0 pr-2 overflow-hidden"
                      >
                        <div className="space-y-1">
                          <div className="font-semibold text-base break-words">{group.name}</div>
                          <div className="text-sm text-muted-foreground line-clamp-2 break-words">
                            {group.specialityName}
                          </div>
                          <div className="text-xs text-muted-foreground line-clamp-1 break-words">
                            {group.facultyName}
                            <span className="hidden sm:inline"> • {group.course} курс</span>
                          </div>
                        </div>
                      </button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 w-8 p-0 shrink-0 flex-shrink-0 mt-0.5"
                        disabled={favoritesLoading}
                        onClick={async (e) => {
                          e.stopPropagation();
                          e.preventDefault();
                          try {
                            await toggleFavorite(group.name);
                          } catch (error) {
                            console.error('Failed to toggle favorite:', error);
                          }
                        }}
                        title={favorite ? 'Удалить из избранного' : 'Добавить в избранное'}
                      >
                        <Star
                          className={`h-4 w-4 ${
                            favorite ? 'fill-yellow-400 text-yellow-400' : 'text-muted-foreground'
                          }`}
                        />
                      </Button>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

