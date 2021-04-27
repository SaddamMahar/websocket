package service

import (
	"messaging-service/model"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"github.com/jinzhu/gorm"
)

func SendingPrivateMessage(message model.Message, onlineUsers map[int64]bool, clients map[*websocket.Conn]int64) {
	saved := true
	for client, id := range clients {
		if id == message.RecipientId || id == message.AuthorId {
			if message.Attachment != "" {
				data, err := ioutil.ReadFile(message.Attachment)
				if err != nil {
					log.Println("error in reading file")
					log.Println(err)
				}
				err = client.WriteMessage(2, data)
				if err != nil {
					log.Println("error in sending file to ", message.RecipientId)
				}
			}
			if saved {
				var singleResp model.MessageResponse
				messgeB, _ := json.Marshal(message)
				err := json.Unmarshal(messgeB, &singleResp)
				if onlineUsers[message.AuthorId] {
					singleResp.Status = true
				}
				err = client.WriteJSON(singleResp)
				if err != nil {
					client.Close()
					delete(clients, client)
					delete(onlineUsers, id)
				}
				saved = false
			}
		}
	}
}

func SendingGroupMessage(message model.Message, clients map[*websocket.Conn]int64) {
	var resp model.MessageResponse
	messgeB, _ := json.Marshal(message)
	json.Unmarshal(messgeB, &resp)

	resp.ActiveUsers = len(clients)

	for client := range clients {
		if message.Attachment != "" {
			data, _ := ioutil.ReadFile(message.Attachment)
			client.WriteMessage(2, data)
		}
		err := client.WriteJSON(resp)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func SendSearchResults(userId int64, search model.Search, onlineUsers map[int64]bool, clients map[*websocket.Conn]int64, db *gorm.DB) {
	if search.Resource == "course" {
		courses := SearchCourse(search, db)
		var resp []model.Course
		for _, v := range courses {
			var course model.Course
			var stdIds []int64
			var instructorIds []int64
			dbCourseBytes, _ := json.Marshal(v)
			json.Unmarshal(dbCourseBytes, &course)
			for _, k := range v.StudentsIDs {
				stdIds = append(stdIds, k.StudentID)
			}
			for _, k := range v.InstructorsIDs {
				instructorIds = append(instructorIds, k.InstructorID)
			}
			course.InstructorsIDs = instructorIds
			course.StudentsIDs = stdIds
			resp = append(resp, course)
		}
		for client, id := range clients {
			if id == userId {
				err := client.WriteJSON(resp)
				if err != nil {
					client.Close()
					delete(clients, client)
					delete(onlineUsers, id)
				}
				break
			}
		}
	} else {
		profiles := SearchProfile(search, db)
		for client, id := range clients {
			if id == userId {
				err := client.WriteJSON(profiles)
				if err != nil {
					client.Close()
					delete(clients, client)
					delete(onlineUsers, id)
				}
				break
			}
		}
	}
}
