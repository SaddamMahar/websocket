package database

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"messaging-service/model"
	"github.com/jinzhu/gorm"
	"time"
)

func SaveMessage(message model.Message, db *gorm.DB) model.Message {
	db.Save(&message)
	return message
}

func UpdateMessage(message model.Message, db *gorm.DB) model.Message {
	var dbMessage model.Message
	db.First(&dbMessage, message.ID)
	dbMessage.Attachment = message.Attachment
	dbMessage.Message = message.Message
	dbMessage.UpdatedAt= time.Now()
	db.Save(&dbMessage)
	return dbMessage
}

func GetMessage(id int64, db *gorm.DB) model.Message {
	var msg model.Message
	db.First(&msg, id)
	return msg
}

func GetMessages(db *gorm.DB) []model.Message {
	var msgs []model.Message
	db.Find(&msgs)
	return msgs
}

func GetMessagesByCourseId(courseId int64, db *gorm.DB) []model.Message {
	var msgs []model.Message
	db.Where("course_id = ?", courseId).Find(&msgs)
	return msgs
}

func GetMessagesByUserId(userId int64, db *gorm.DB) []model.Message {
	var msgs []model.Message
	db.Where("author_id = ?", userId).Find(&msgs).Order("created_at DESC")
	return msgs
}

func SaveCourse(course model.DBCourse, db *gorm.DB) model.DBCourse {
	db.Create(&course)
	return course
}

func UpdateStudentIds(courseId int64, stds []model.StudentsIDs, db *gorm.DB) model.DBCourse {
	var course model.DBCourse
	db.First(&course, courseId)
	db.Model(&course).Association("StudentsIDs").Append(stds)
	return course
}

func GetCourse(id int64, db *gorm.DB) model.DBCourse {
	var course model.DBCourse
	var instructors []model.InstructorsIDs
	var students []model.StudentsIDs
	db.First(&course, id)
	db.Model(&course).Related(&instructors, "InstructorsIDs")
	db.Model(&course).Related(&students, "StudentsIDs")

	course.InstructorsIDs = instructors
	course.StudentsIDs = students
	return course
}

func GetCoursesByField(fieldName, value string, db *gorm.DB) []model.DBCourse {
	var course []model.DBCourse
	var instructors []model.InstructorsIDs
	var students []model.StudentsIDs

	db.Where(fieldName+ " LIKE ?", "%"+value+"%").Find(&course)
	for i,v := range course {
		db.Model(&v).Related(&instructors)
		db.Model(&v).Related(&students)

		v.InstructorsIDs = instructors
		v.StudentsIDs = students
		course[i] = v
	}

	return course
}

func GetCoursesByStudent(studentId int64, db *gorm.DB) []model.StudentsIDs {
	var students []model.StudentsIDs
	db.Find(&students,studentId)
	return students
}

func SaveProfile(profile model.Profile, db *gorm.DB) model.Profile {
	db.Create(&profile)
	return profile
}

func GetProfileById(id int64, db *gorm.DB) model.Profile {
	var profile model.Profile
	db.First(&profile, id)
	return profile
}

func GetProfileByName(name string, db *gorm.DB) []model.Profile {
	var profiles []model.Profile
	db.Where("name LIKE ?", name).Find(&profiles).Order("updated_at DESC")
	return profiles
}