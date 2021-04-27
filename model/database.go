package model

import (
	"time"
	"github.com/jinzhu/gorm"
)

type Message struct {
	ID          int64  `json:"id"`
	AuthorId    int64  `json:"author_id"`
	CourseId    int64  `json:"course_id"`
	RecipientId int64  `json:"recipient_id"`
	Message     string `json:"message"`
	Attachment  string `json:"attachment"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Profile struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Avatar     string `json:"avatarUrl"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DBCourse struct {
	ID             int64 `gorm:"primaryKey"`
	Name           string           `json:"name"`
	CardPhotoURL   string           `json:"card_photo_url"`
	AuthorID       int64            `json:"author_id"`
	InstructorsIDs []InstructorsIDs `json:"instructors_ids" gorm:"ForeignKey:CourseID"`
	StudentsIDs    []StudentsIDs    `json:"students_ids" gorm:"ForeignKey:CourseID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StudentsIDs struct {
	gorm.Model
	CourseId  int64 `json:"course_id"`
	StudentID int64 `json:"student_id"`
}

type InstructorsIDs struct {
	gorm.Model
	CourseId     int64 `json:"course_id"`
	InstructorID int64 `json:"instructor_id"`
}

