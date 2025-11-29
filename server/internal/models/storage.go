package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoredSchedule struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupNumber    string             `bson:"group_number" json:"group_number"`
	EmployeeURLID  string             `bson:"employee_url_id,omitempty" json:"employee_url_id,omitempty"`
	ScheduleData   ScheduleResponse   `bson:"schedule_data" json:"schedule_data"`
	LastUpdateDate string             `bson:"last_update_date" json:"last_update_date"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type StoredGroup struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	BSUIRID        int                  `bson:"bsuir_id" json:"bsuir_id"`
	GroupData      StudentGroupListItem `bson:"group_data" json:"group_data"`
	LastUpdateDate string               `bson:"last_update_date" json:"last_update_date"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updated_at"`
}

type StoredEmployee struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BSUIRID        int                `bson:"bsuir_id" json:"bsuir_id"`
	URLID          string             `bson:"url_id" json:"url_id"`
	EmployeeData   EmployeeListItem   `bson:"employee_data" json:"employee_data"`
	LastUpdateDate string             `bson:"last_update_date" json:"last_update_date"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type StoredFaculty struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BSUIRID     int                `bson:"bsuir_id" json:"bsuir_id"`
	FacultyData Faculty            `bson:"faculty_data" json:"faculty_data"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type StoredDepartment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BSUIRID        int                `bson:"bsuir_id" json:"bsuir_id"`
	DepartmentData Department         `bson:"department_data" json:"department_data"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type StoredSpeciality struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BSUIRID        int                `bson:"bsuir_id" json:"bsuir_id"`
	SpecialityData Speciality         `bson:"speciality_data" json:"speciality_data"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type StoredAuditory struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BSUIRID      int                `bson:"bsuir_id" json:"bsuir_id"`
	AuditoryData Auditory           `bson:"auditory_data" json:"auditory_data"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
