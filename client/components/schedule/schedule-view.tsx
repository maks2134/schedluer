'use client';

import { useState, useMemo } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { RefreshCw } from 'lucide-react';
import type { ScheduleResponse, Schedule } from '@/types/api';

interface ScheduleViewProps {
  schedule: ScheduleResponse | null;
  loading: boolean;
  onRefresh?: () => void;
  refreshing?: boolean;
}

const weekDays = ['Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота'];

const lessonTypeColors: Record<string, string> = {
  'ЛК': 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
  'ЛР': 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  'ПЗ': 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300',
  'Экзамен': 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
  'Консультация': 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
};

// Функция для удаления дубликатов пар (учитывает подгруппы и недели)
function removeDuplicateLessons(lessons: Schedule[]): Schedule[] {
  const seen = new Set<string>();
  const unique: Schedule[] = [];

  for (const lesson of lessons) {
    // Создаем уникальный ключ: время + предмет + тип + аудитория + подгруппа + недели
    const weekKey = lesson.weekNumber && Array.isArray(lesson.weekNumber) 
      ? lesson.weekNumber.sort((a, b) => a - b).join(',')
      : '';
    const key = `${lesson.startLessonTime}-${lesson.endLessonTime}-${lesson.subject}-${lesson.lessonTypeAbbrev}-${(lesson.auditories || []).join(',')}-${lesson.numSubgroup}-${weekKey}`;
    
    if (!seen.has(key)) {
      seen.add(key);
      unique.push(lesson);
    }
  }

  return unique;
}

// Функция для фильтрации по подгруппе
function filterBySubgroup(lessons: Schedule[], subgroup: number): Schedule[] {
  if (subgroup === 0) {
    return lessons; // Показываем все
  }
  return lessons.filter((lesson) => lesson.numSubgroup === subgroup || lesson.numSubgroup === 0);
}

export function ScheduleView({ schedule, loading, onRefresh, refreshing }: ScheduleViewProps) {
  const [selectedSubgroup, setSelectedSubgroup] = useState<number>(0); // 0 = все, 1 = подгруппа 1, 2 = подгруппа 2

  // Получаем доступные подгруппы из расписания
  const availableSubgroups = useMemo(() => {
    if (!schedule?.schedules) return [0];
    const subgroups = new Set<number>();
    Object.values(schedule.schedules).forEach((daySchedule) => {
      if (Array.isArray(daySchedule)) {
        daySchedule.forEach((lesson) => {
          if (lesson.numSubgroup > 0) {
            subgroups.add(lesson.numSubgroup);
          }
        });
      }
    });
    const result = [0, ...Array.from(subgroups).sort()];
    return result;
  }, [schedule]);

  if (loading) {
    return (
      <div className="space-y-4">
        {weekDays.map((day) => (
          <Card key={day}>
            <CardHeader>
              <Skeleton className="h-6 w-32" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-32 w-full" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (!schedule) {
    return (
      <Card>
        <CardContent className="py-8 text-center text-muted-foreground">
          Выберите группу или преподавателя для просмотра расписания
        </CardContent>
      </Card>
    );
  }

  const groupInfo = schedule.studentGroupDto;
  const employeeInfo = schedule.employeeDto;

  return (
    <div className="space-y-6">
      {/* Header */}
      <Card>
        <CardHeader>
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div className="min-w-0 flex-1">
              <CardTitle className="text-lg sm:text-xl truncate">
                {groupInfo ? `Группа ${groupInfo.name}` : employeeInfo?.fio || 'Расписание'}
              </CardTitle>
              <CardDescription className="text-xs sm:text-sm">
                {groupInfo && (
                  <>
                    <span className="hidden sm:inline">{groupInfo.specialityName} • </span>
                    {groupInfo.course} курс • {groupInfo.facultyAbbrev}
                  </>
                )}
                {employeeInfo && (
                  <>
                    {employeeInfo.rank} • {employeeInfo.degree}
                  </>
                )}
              </CardDescription>
            </div>
            {onRefresh && (
              <Button
                variant="outline"
                size="sm"
                onClick={onRefresh}
                disabled={refreshing}
                className="shrink-0"
              >
                <RefreshCw className={`h-4 w-4 sm:mr-2 ${refreshing ? 'animate-spin' : ''}`} />
                <span className="hidden sm:inline">Обновить</span>
              </Button>
            )}
          </div>
        </CardHeader>
        {(schedule.startDate || schedule.endDate) && (
          <CardContent>
            <div className="text-sm text-muted-foreground">
              Период: {schedule.startDate} - {schedule.endDate}
            </div>
          </CardContent>
        )}
      </Card>

      {/* Schedule with Tabs */}
      <Tabs defaultValue="lessons" className="w-full">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="lessons">Пары</TabsTrigger>
          <TabsTrigger value="exams">Экзамены</TabsTrigger>
        </TabsList>
        
        <TabsContent value="lessons" className="mt-4">
          {schedule.schedules && Object.keys(schedule.schedules).length > 0 ? (
            <div className="space-y-4">
              {/* Фильтр подгрупп */}
              {availableSubgroups.length > 1 && (
                <Card>
                  <CardContent className="pt-4 sm:pt-6">
                    <div className="flex flex-col sm:flex-row sm:items-center gap-2">
                      <span className="text-sm font-medium">Подгруппа:</span>
                      <Select
                        value={selectedSubgroup.toString()}
                        onValueChange={(value) => setSelectedSubgroup(Number(value))}
                      >
                        <SelectTrigger className="w-full sm:w-[180px]">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="0">Все</SelectItem>
                          {availableSubgroups.filter(s => s > 0).map((sg) => (
                            <SelectItem key={sg} value={sg.toString()}>
                              Подгруппа {sg}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  </CardContent>
                </Card>
              )}

              {weekDays.map((day) => {
                const daySchedule = schedule.schedules?.[day] || [];
                if (!Array.isArray(daySchedule) || daySchedule.length === 0) return null;

                // Дедупликация пар по уникальным характеристикам
                const uniqueLessons = removeDuplicateLessons(daySchedule);
                // Фильтрация по выбранной подгруппе
                const filteredLessons = filterBySubgroup(uniqueLessons, selectedSubgroup);

                if (filteredLessons.length === 0) return null;

                return (
                  <Card key={day}>
                    <CardHeader>
                      <CardTitle className="text-lg">{day}</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <div className="space-y-3">
                        {filteredLessons.map((lesson: Schedule, idx: number) => (
                          <LessonCard key={`${day}-${idx}-${lesson.startLessonTime}-${lesson.subject}-${lesson.numSubgroup}`} lesson={lesson} />
                        ))}
                      </div>
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          ) : (
            <Card>
              <CardContent className="py-8 text-center text-muted-foreground">
                Расписание пар отсутствует
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="exams" className="mt-4">
          {schedule.exams && Array.isArray(schedule.exams) && schedule.exams.length > 0 ? (
            <Card>
              <CardHeader>
                <CardTitle>Экзамены</CardTitle>
                {(schedule.startExamsDate || schedule.endExamsDate) && (
                  <CardDescription>
                    Период: {schedule.startExamsDate} - {schedule.endExamsDate}
                  </CardDescription>
                )}
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {schedule.exams.map((exam: Schedule, idx: number) => (
                    <LessonCard key={idx} lesson={exam} />
                  ))}
                </div>
              </CardContent>
            </Card>
          ) : (
            <Card>
              <CardContent className="py-8 text-center text-muted-foreground">
                Экзамены отсутствуют
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}

function LessonCard({ lesson }: { lesson: Schedule }) {
  return (
    <div className="flex flex-col sm:flex-row items-start gap-3 sm:gap-4 p-3 border rounded-lg hover:bg-accent/50 transition-colors">
      <div className="flex-shrink-0 w-full sm:w-20 text-sm font-medium">
        {lesson.startLessonTime} - {lesson.endLessonTime}
      </div>
      <div className="flex-1 min-w-0 w-full">
        <div className="flex items-start justify-between gap-2 mb-1">
          <div>
            <div className="font-semibold">{lesson.subjectFullName || lesson.subject}</div>
            {lesson.note && (
              <div className="text-sm text-muted-foreground">{lesson.note}</div>
            )}
          </div>
          <Badge
            variant="outline"
            className={lessonTypeColors[lesson.lessonTypeAbbrev] || ''}
          >
            {lesson.lessonTypeAbbrev}
          </Badge>
        </div>
        <div className="flex flex-wrap gap-2 text-sm text-muted-foreground">
          {lesson.auditories && Array.isArray(lesson.auditories) && lesson.auditories.length > 0 && (
            <span>📍 {lesson.auditories.join(', ')}</span>
          )}
          {lesson.employees && Array.isArray(lesson.employees) && lesson.employees.length > 0 && (
            <span>
              👤 {lesson.employees.map((e) => `${e.lastName} ${e.firstName[0]}. ${e.middleName[0]}.`).join(', ')}
            </span>
          )}
          {lesson.numSubgroup > 0 && (
            <span>Подгруппа {lesson.numSubgroup}</span>
          )}
          {lesson.weekNumber && Array.isArray(lesson.weekNumber) && lesson.weekNumber.length > 0 && (
            <span>Недели: {lesson.weekNumber.join(', ')}</span>
          )}
        </div>
        {lesson.dateLesson && (
          <div className="text-xs text-muted-foreground mt-1">
            Дата: {lesson.dateLesson}
          </div>
        )}
      </div>
    </div>
  );
}

