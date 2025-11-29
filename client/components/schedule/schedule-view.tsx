'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
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

export function ScheduleView({ schedule, loading, onRefresh, refreshing }: ScheduleViewProps) {
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
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>
                {groupInfo ? `Группа ${groupInfo.name}` : employeeInfo?.fio || 'Расписание'}
              </CardTitle>
              <CardDescription>
                {groupInfo && (
                  <>
                    {groupInfo.specialityName} • {groupInfo.course} курс • {groupInfo.facultyAbbrev}
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
              >
                <RefreshCw className={`h-4 w-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
                Обновить
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

      {/* Schedule by days */}
      {schedule.schedules && Object.keys(schedule.schedules).length > 0 ? (
        <div className="space-y-4">
          {weekDays.map((day) => {
            const daySchedule = schedule.schedules?.[day] || [];
            if (!Array.isArray(daySchedule) || daySchedule.length === 0) return null;

            return (
              <Card key={day}>
                <CardHeader>
                  <CardTitle className="text-lg">{day}</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    {daySchedule.map((lesson: Schedule, idx: number) => (
                      <LessonCard key={idx} lesson={lesson} />
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
            Расписание отсутствует
          </CardContent>
        </Card>
      )}

      {/* Exams */}
      {schedule.exams && Array.isArray(schedule.exams) && schedule.exams.length > 0 && (
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
      )}
    </div>
  );
}

function LessonCard({ lesson }: { lesson: Schedule }) {
  return (
    <div className="flex items-start gap-4 p-3 border rounded-lg hover:bg-accent/50 transition-colors">
      <div className="flex-shrink-0 w-20 text-sm font-medium">
        {lesson.startLessonTime} - {lesson.endLessonTime}
      </div>
      <div className="flex-1 min-w-0">
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

