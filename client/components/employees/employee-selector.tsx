'use client';

import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Search, RefreshCw } from 'lucide-react';
import type { EmployeeListItem } from '@/types/api';

interface EmployeeSelectorProps {
  employees: EmployeeListItem[];
  loading: boolean;
  selectedEmployee: string | null;
  onSelectEmployee: (urlId: string) => void;
  onRefresh?: () => void;
  refreshing?: boolean;
}

export function EmployeeSelector({
  employees,
  loading,
  selectedEmployee,
  onSelectEmployee,
  onRefresh,
  refreshing,
}: EmployeeSelectorProps) {
  const [searchQuery, setSearchQuery] = useState('');

  const filteredEmployees = employees.filter((employee) =>
    employee.fio.toLowerCase().includes(searchQuery.toLowerCase()) ||
    employee.lastName.toLowerCase().includes(searchQuery.toLowerCase()) ||
    employee.firstName.toLowerCase().includes(searchQuery.toLowerCase()) ||
    employee.rank.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const getInitials = (employee: EmployeeListItem) => {
    return `${employee.firstName[0]}${employee.lastName[0]}`.toUpperCase();
  };

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
            <CardTitle>Преподаватели</CardTitle>
            <CardDescription>
              Найдено: {filteredEmployees.length} из {employees.length}
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
              placeholder="Поиск по ФИО, должности..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9"
            />
          </div>
          <div className="max-h-[600px] overflow-y-auto space-y-2">
            {filteredEmployees.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                Преподаватели не найдены
              </div>
            ) : (
              filteredEmployees.map((employee) => (
                <button
                  key={employee.id}
                  onClick={() => onSelectEmployee(employee.urlId)}
                  className={`w-full text-left p-3 rounded-lg border transition-all hover:bg-accent ${
                    selectedEmployee === employee.urlId
                      ? 'border-primary bg-accent'
                      : 'border-border'
                  }`}
                >
                  <div className="flex items-center gap-3">
                    <Avatar>
                      <AvatarImage src={employee.photoLink} alt={employee.fio} />
                      <AvatarFallback>{getInitials(employee)}</AvatarFallback>
                    </Avatar>
                    <div className="flex-1 min-w-0">
                      <div className="font-semibold">{employee.fio}</div>
                      <div className="text-sm text-muted-foreground">
                        {employee.rank}
                      </div>
                      {employee.degree && (
                        <div className="text-xs text-muted-foreground mt-1">
                          {employee.degree}
                        </div>
                      )}
                    </div>
                    {employee.academicDepartment && employee.academicDepartment.length > 0 && (
                      <Badge variant="outline" className="text-xs">
                        {employee.academicDepartment[0]}
                      </Badge>
                    )}
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

