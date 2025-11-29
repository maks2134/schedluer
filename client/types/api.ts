// API Types - соответствуют моделям из backend

export interface ScheduleResponse {
  employeeDto: EmployeeDto | null;
  studentGroupDto: StudentGroupDto | null;
  schedules: Record<string, Schedule[]>;
  exams: Schedule[];
  startDate: string | null;
  endDate: string | null;
  startExamsDate: string | null;
  endExamsDate: string | null;
}

export interface EmployeeDto {
  firstName: string;
  lastName: string;
  middleName: string;
  degree: string;
  degreeAbbrev: string;
  rank: string;
  photoLink: string;
  calendarId: string;
  id: number;
  urlId: string;
  email: string;
  jobPositions: string[];
  fio?: string;
  academicDepartment?: string[];
}

export interface StudentGroupDto {
  name: string;
  facultyId: number;
  facultyAbbrev: string;
  specialityDepartmentEducationFormId: number;
  specialityName: string;
  specialityAbbrev: string;
  course: number;
  id: number;
  calendarId: string;
  educationDegree: number;
  facultyName?: string;
}

export interface Schedule {
  weekNumber: number[];
  studentGroups: StudentGroup[];
  numSubgroup: number;
  auditories: string[];
  startLessonTime: string;
  endLessonTime: string;
  subject: string;
  subjectFullName: string;
  note: string | null;
  lessonTypeAbbrev: string;
  dateLesson: string | null;
  startLessonDate: string | null;
  endLessonDate: string | null;
  announcement: boolean;
  split: boolean;
  employees: EmployeeDto[];
}

export interface StudentGroup {
  specialityName: string;
  specialityCode: string;
  numberOfStudents: number;
  name: string;
  educationDegree: number;
}

export interface StudentGroupListItem {
  name: string;
  facultyId: number;
  facultyName: string;
  specialityDepartmentEducationFormId: number;
  specialityName: string;
  course: number;
  id: number;
  calendarId: string;
}

export interface EmployeeListItem {
  firstName: string;
  lastName: string;
  middleName: string;
  degree: string;
  rank: string;
  photoLink: string;
  calendarId: string;
  academicDepartment: string[];
  id: number;
  urlId: string;
  fio: string;
}

export interface ApiError {
  error: string;
}

export interface RefreshResponse {
  message: string;
}

