'use client';

import { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { GroupSelector } from '@/components/groups/group-selector';
import { EmployeeSelector } from '@/components/employees/employee-selector';
import { ScheduleView } from '@/components/schedule/schedule-view';
import { useGroups } from '@/lib/hooks/use-groups';
import { useEmployees } from '@/lib/hooks/use-employees';
import { useGroupSchedule } from '@/lib/hooks/use-schedule';
import { useEmployeeSchedule } from '@/lib/hooks/use-schedule';

export default function Home() {
  const [selectedGroup, setSelectedGroup] = useState<string | null>(null);
  const [selectedEmployee, setSelectedEmployee] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'groups' | 'employees'>('groups');

  const { data: groups, loading: groupsLoading, refresh: refreshGroups, error: groupsError } = useGroups();
  const { data: employees, loading: employeesLoading, refresh: refreshEmployees, error: employeesError } = useEmployees();
  
  const { 
    data: groupSchedule, 
    loading: groupScheduleLoading, 
    refresh: refreshGroupSchedule,
    error: groupScheduleError 
  } = useGroupSchedule(selectedGroup);
  
  const { 
    data: employeeSchedule, 
    loading: employeeScheduleLoading, 
    refresh: refreshEmployeeSchedule,
    error: employeeScheduleError 
  } = useEmployeeSchedule(selectedEmployee);

  const handleGroupSelect = (groupNumber: string) => {
    setSelectedGroup(groupNumber);
    setSelectedEmployee(null);
    setActiveTab('groups');
  };

  const handleEmployeeSelect = (urlId: string) => {
    setSelectedEmployee(urlId);
    setSelectedGroup(null);
    setActiveTab('employees');
  };

  const currentSchedule = selectedGroup ? groupSchedule : employeeSchedule;
  const scheduleLoading = selectedGroup ? groupScheduleLoading : employeeScheduleLoading;
  const scheduleError = selectedGroup ? groupScheduleError : employeeScheduleError;
  const handleRefresh = selectedGroup 
    ? () => refreshGroupSchedule() 
    : () => refreshEmployeeSchedule();

  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto py-8 px-4">
        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Расписание БГУИР</h1>
          <p className="text-muted-foreground">
            Просмотр расписания групп и преподавателей
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 sm:gap-6">
          {/* Left sidebar - Groups/Employees */}
          <div className="lg:col-span-1 order-2 lg:order-1">
            <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as 'groups' | 'employees')}>
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="groups">Группы</TabsTrigger>
                <TabsTrigger value="employees">Преподаватели</TabsTrigger>
              </TabsList>
              <TabsContent value="groups" className="mt-4">
                <GroupSelector
                  groups={groups}
                  loading={groupsLoading}
                  selectedGroup={selectedGroup}
                  onSelectGroup={handleGroupSelect}
                  onRefresh={refreshGroups}
                />
                {groupsError && (
                  <div className="mt-4 p-3 bg-destructive/10 text-destructive rounded-lg text-sm">
                    {groupsError}
                  </div>
                )}
              </TabsContent>
              <TabsContent value="employees" className="mt-4">
                <EmployeeSelector
                  employees={employees}
                  loading={employeesLoading}
                  selectedEmployee={selectedEmployee}
                  onSelectEmployee={handleEmployeeSelect}
                  onRefresh={refreshEmployees}
            />
                {employeesError && (
                  <div className="mt-4 p-3 bg-destructive/10 text-destructive rounded-lg text-sm">
                    {employeesError}
                  </div>
                )}
              </TabsContent>
            </Tabs>
          </div>

          {/* Right side - Schedule */}
          <div className="lg:col-span-2 order-1 lg:order-2">
            {scheduleError && (
              <div className="mb-4 p-4 bg-destructive/10 text-destructive rounded-lg">
                Ошибка: {scheduleError}
              </div>
            )}
            <ScheduleView
              schedule={currentSchedule}
              loading={scheduleLoading}
              onRefresh={handleRefresh}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
