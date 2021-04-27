package service

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"messaging-service/model"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/websocket"
	"messaging-service/database"
	"github.com/nats-io/nats.go"
	"os"
	"strconv"
	"log"
)

func SaveMessage(message model.Message, conn *nats.Conn, db *gorm.DB) model.Message {
	message = database.SaveMessage(message, db)
	return message
}

func SaveCourse(course model.Course, db *gorm.DB) {
	var dbCourse model.DBCourse
	var studentsIDs []model.StudentsIDs
	var instructorsIDs []model.InstructorsIDs
	if len(course.StudentsIDs) > 0 {
		for _, v := range course.StudentsIDs {
			var stdId model.StudentsIDs
			stdId.CourseId = course.ID
			stdId.StudentID = v
			studentsIDs = append(studentsIDs, stdId)
		}
	}
	if len(course.InstructorsIDs) > 0 {
		for _, v := range course.InstructorsIDs {
			var instructorID model.InstructorsIDs
			instructorID.CourseId = course.ID
			instructorID.InstructorID = v
			instructorsIDs = append(instructorsIDs, instructorID)
		}
	}
	dbCourse.StudentsIDs = studentsIDs
	dbCourse.InstructorsIDs = instructorsIDs

	dbCourse.ID = course.ID
	dbCourse.Name = course.Name
	dbCourse.AuthorID = course.AuthorID
	dbCourse.CardPhotoURL = course.CardPhotoURL

	database.SaveCourse(dbCourse, db)
}

func SendAllCourseConversation(userId int64, client *websocket.Conn, db *gorm.DB) {
	sutudentWithCourseIds := database.GetCoursesByStudent(userId, db)

	for _, v := range sutudentWithCourseIds {
		messages := database.GetMessagesByCourseId(v.CourseId, db)
		err := client.WriteJSON(messages)
		if err != nil {
			client.Close()
		}
	}
}

func SendAllConversations(userId int64, client *websocket.Conn, db *gorm.DB) {
	messages := database.GetMessagesByUserId(userId, db)
	err := client.WriteJSON(messages)
	if err != nil {
		client.Close()
	}
}

func GetFileFromWS(data []byte, userId int64, db *gorm.DB) model.Message {
	msgs := database.GetMessagesByUserId(userId, db)
	path, err := os.Getwd()
	if err != nil {
		log.Println("Unable to get current path, "+err.Error())
	}
	path = path + "\\uploads"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	path = path + "\\" + strconv.FormatInt(userId, 10)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	var message model.Message
	for _, v := range msgs {
		if v.Attachment != "" {
			file, err := os.Create(path + "\\" + v.Attachment)
			if err != nil {
				log.Println("Unable to create file, "+err.Error())
			}
			defer file.Close()
			file.Write(data)
			v.Attachment = path + "\\" + v.Attachment
			message = database.UpdateMessage(v, db)
			break
		}
	}
	return message
}

func UpdateStudentsForCourse(courseStudents model.CourseStudents, db *gorm.DB) {
	var stds []model.StudentsIDs
	for _, v := range courseStudents.StudentIDs {
		var std model.StudentsIDs
		std.CourseId = courseStudents.CourseId
		std.StudentID = v
		stds = append(stds, std)
	}
	database.UpdateStudentIds(courseStudents.CourseId, stds, db)
}

func GetCourseInfo(id int64, db *gorm.DB) model.DBCourse {
	return database.GetCourse(id, db)
}

func SaveProfileToDB(profile model.Profile, db *gorm.DB) model.Profile {
	return database.SaveProfile(profile, db)
}

func GetProfile(id int64, db *gorm.DB) model.Profile {
	return database.GetProfileById(id, db)
}