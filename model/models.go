package model

import (
	"time"
)

type Search struct {
	Name     string `json:"name"`
	Resource string `json:"resource"`
}

type Course struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	CardPhotoURL   string  `json:"card_photo_url"`
	AuthorID       int64   `json:"author_id"`
	InstructorsIDs []int64 `json:"instructors_ids"`
	StudentsIDs    []int64 `json:"students_ids"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProfileWrapper struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Profile Profile
}

type MessageResponse struct {
	ID            int64   `json:"id"`
	AuthorID      int64   `json:"author_id"`
	AuthorProfile Profile `json:"profile"`
	CourseID      int64   `json:"course_id"`
	RecipientID   int64   `json:"recipient_id"`
	Message       string  `json:"message"`

	Status      bool `json:"status"`
	ActiveUsers int  `json:"active_users"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchResponse struct {
	Profiles []Profile `json:"profiles"`
}

type CourseStudents struct {
	CourseId   int64   `json:"course_id"`
	StudentIDs []int64 `json:"student_ids"`
}

type CourseInstructors struct {
	CourseId      int64   `json:"course_id"`
	InstructorIDs []int64 `json:"instructor_ids"`
}
