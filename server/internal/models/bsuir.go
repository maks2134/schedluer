package models

type ScheduleResponse struct {
	EmployeeDto     *EmployeeDto          `json:"employeeDto"`
	StudentGroupDto *StudentGroupDto      `json:"studentGroupDto"`
	Schedules       map[string][]Schedule `json:"schedules"`
	Exams           []Schedule            `json:"exams"`
	StartDate       string                `json:"startDate"`
	EndDate         string                `json:"endDate"`
	StartExamsDate  string                `json:"startExamsDate"`
	EndExamsDate    string                `json:"endExamsDate"`
}

type EmployeeDto struct {
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	MiddleName   string   `json:"middleName"`
	Degree       string   `json:"degree"`
	DegreeAbbrev string   `json:"degreeAbbrev"`
	Rank         string   `json:"rank"`
	PhotoLink    string   `json:"photoLink"`
	CalendarID   string   `json:"calendarId"`
	ID           int      `json:"id"`
	URLID        string   `json:"urlId"`
	Email        string   `json:"email"`
	JobPositions []string `json:"jobPositions"`
	FIO          string   `json:"fio,omitempty"`
	AcademicDept []string `json:"academicDepartment,omitempty"`
}

type StudentGroupDto struct {
	Name                                string `json:"name"`
	FacultyID                           int    `json:"facultyId"`
	FacultyAbbrev                       string `json:"facultyAbbrev"`
	SpecialityDepartmentEducationFormID int    `json:"specialityDepartmentEducationFormId"`
	SpecialityName                      string `json:"specialityName"`
	SpecialityAbbrev                    string `json:"specialityAbbrev"`
	Course                              int    `json:"course"`
	ID                                  int    `json:"id"`
	CalendarID                          string `json:"calendarId"`
	EducationDegree                     int    `json:"educationDegree"`
	FacultyName                         string `json:"facultyName,omitempty"`
}

type Schedule struct {
	WeekNumber       []int          `json:"weekNumber"`
	StudentGroups    []StudentGroup `json:"studentGroups"`
	NumSubgroup      int            `json:"numSubgroup"`
	Auditories       []string       `json:"auditories"`
	StartLessonTime  string         `json:"startLessonTime"`
	EndLessonTime    string         `json:"endLessonTime"`
	Subject          string         `json:"subject"`
	SubjectFullName  string         `json:"subjectFullName"`
	Note             string         `json:"note"`
	LessonTypeAbbrev string         `json:"lessonTypeAbbrev"`
	DateLesson       string         `json:"dateLesson"`
	StartLessonDate  string         `json:"startLessonDate"`
	EndLessonDate    string         `json:"endLessonDate"`
	Announcement     bool           `json:"announcement"`
	Split            bool           `json:"split"`
	Employees        []EmployeeDto  `json:"employees"`
}

type StudentGroup struct {
	SpecialityName   string `json:"specialityName"`
	SpecialityCode   string `json:"specialityCode"`
	NumberOfStudents int    `json:"numberOfStudents"`
	Name             string `json:"name"`
	EducationDegree  int    `json:"educationDegree"`
}

type StudentGroupListItem struct {
	Name                                string `json:"name"`
	FacultyID                           int    `json:"facultyId"`
	FacultyName                         string `json:"facultyName"`
	SpecialityDepartmentEducationFormID int    `json:"specialityDepartmentEducationFormId"`
	SpecialityName                      string `json:"specialityName"`
	Course                              int    `json:"course"`
	ID                                  int    `json:"id"`
	CalendarID                          string `json:"calendarId"`
}

type EmployeeListItem struct {
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	MiddleName   string   `json:"middleName"`
	Degree       string   `json:"degree"`
	Rank         string   `json:"rank"`
	PhotoLink    string   `json:"photoLink"`
	CalendarID   string   `json:"calendarId"`
	AcademicDept []string `json:"academicDepartment"`
	ID           int      `json:"id"`
	URLID        string   `json:"urlId"`
	FIO          string   `json:"fio"`
}

type Faculty struct {
	Name   string `json:"name"`
	Abbrev string `json:"abbrev"`
	ID     int    `json:"id"`
}

type Department struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Abbrev string `json:"abbrev"`
}

type Speciality struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	Abbrev        string          `json:"abbrev"`
	EducationForm []EducationForm `json:"educationForm"`
	FacultyID     int             `json:"facultyId"`
	Code          string          `json:"code"`
}

type EducationForm struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Announcement struct {
	ID                  int                 `json:"id"`
	Employee            string              `json:"employee"`
	Content             string              `json:"content"`
	Date                string              `json:"date"`
	EmployeeDepartments []string            `json:"employeeDepartments"`
	StudentGroups       []AnnouncementGroup `json:"studentGroups"`
}

type AnnouncementGroup struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Auditory struct {
	ID             int                `json:"id"`
	Name           string             `json:"name"`
	Note           string             `json:"note"`
	Capacity       *int               `json:"capacity"`
	AuditoryType   AuditoryType       `json:"auditoryType"`
	BuildingNumber BuildingNumber     `json:"buildingNumber"`
	Department     AuditoryDepartment `json:"department"`
}

type AuditoryType struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Abbrev string `json:"abbrev"`
}

type BuildingNumber struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AuditoryDepartment struct {
	IDDepartment  int    `json:"idDepartment"`
	Abbrev        string `json:"abbrev"`
	Name          string `json:"name"`
	NameAndAbbrev string `json:"nameAndAbbrev"`
}

type LastUpdateDate struct {
	LastUpdateDate string `json:"lastUpdateDate"`
}
