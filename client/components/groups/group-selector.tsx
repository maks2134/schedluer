'use client';

import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Search, RefreshCw } from 'lucide-react';
import type { StudentGroupListItem } from '@/types/api';

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

  const filteredGroups = groups.filter((group) =>
    group.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    group.specialityName.toLowerCase().includes(searchQuery.toLowerCase()) ||
    group.facultyName.toLowerCase().includes(searchQuery.toLowerCase())
  );

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
              Найдено: {filteredGroups.length} из {groups.length}
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
              placeholder="Поиск по номеру, специальности, факультету..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9"
            />
          </div>
          <div className="max-h-[600px] overflow-y-auto space-y-2">
            {filteredGroups.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                Группы не найдены
              </div>
            ) : (
              filteredGroups.map((group) => (
                <button
                  key={group.id}
                  onClick={() => onSelectGroup(group.name)}
                  className={`w-full text-left p-3 rounded-lg border transition-all hover:bg-accent ${
                    selectedGroup === group.name
                      ? 'border-primary bg-accent'
                      : 'border-border'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <div>
                      <div className="font-semibold">{group.name}</div>
                      <div className="text-sm text-muted-foreground">
                        {group.specialityName}
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">
                        {group.facultyName} • {group.course} курс
                      </div>
                    </div>
                    <Badge variant="outline">{group.course} курс</Badge>
                  </div>
                </button>
              ))
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

