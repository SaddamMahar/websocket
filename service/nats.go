package service

import (
	"log"
	"github.com/nats-io/nats.go"
	"messaging-service/model"
	"encoding/json"
	"time"
	"github.com/jinzhu/gorm"
)

func ConnectNATS() *nats.Conn {
	nc, err := nats.Connect(model.NATS_URL)
	defer nc.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to NATS SERVER")
	return nc
}

func UpdateCourseInfo(db *gorm.DB, conn *nats.Conn) {
	var course model.Course
	var err error
	conn.Subscribe("courseChat.UpdateCourseInfo", func(msg *nats.Msg) {
		err = json.Unmarshal(msg.Data, &course)
		if err != nil {
			berr, _ := json.Marshal(err)
			conn.Publish("courseChat.UpdateCourseInfo", berr)
		}
		conn.Publish("courseChat.UpdateCourseInfo", []byte{})
		SaveCourse(course, db)
	})
}

func UpdateCourseStudents(db *gorm.DB, conn *nats.Conn) {
	var courseStudents model.CourseStudents
	var err error
	conn.Subscribe("courseChat.UpdateCourseStudents", func(msg *nats.Msg) {
		err = json.Unmarshal(msg.Data, &courseStudents)
		if err != nil {
			berr, _ := json.Marshal(err)
			conn.Publish("courseChat.UpdateCourseStudents", berr)
		}
		conn.Publish("courseChat.UpdateCourseStudents", []byte{})
		UpdateStudentsForCourse(courseStudents, db)
	})
}

func GetProfileFromNATS(profileId int64, conn *nats.Conn) (model.ProfileWrapper, error) {
	var userProfile model.ProfileWrapper
	byteBody, _ := json.Marshal(profileId)
	msg, err := conn.Request("userProfile.GetProfile", byteBody, 10*time.Second)
	if err != nil {
		return userProfile, err
	}
	err = json.Unmarshal(msg.Data, &userProfile)
	return userProfile, err
}
