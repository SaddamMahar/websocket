package main

import (
	"net/http"
	"log"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"messaging-service/model"
	"messaging-service/service"
	"github.com/nats-io/nats.go"
	"strconv"
	"encoding/json"
	"strings"
)

var db *gorm.DB
var natsConn *nats.Conn

var err error
var activeUsers = make(map[*websocket.Conn]int64)
var onlineUsers = make(map[int64]bool)
var broadcast = make(chan model.Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	db, err = gorm.Open("postgres",
		"host=localhost port=5432 user=postgres dbname=postgres sslmode=disable password=postgres")

	if err != nil {
		panic("failed to connect database")

	}
	defer db.Close()
	db.AutoMigrate(&model.Message{})
	db.AutoMigrate(&model.DBCourse{})
	db.AutoMigrate(&model.Profile{})
	db.AutoMigrate(&model.StudentsIDs{})
	db.AutoMigrate(&model.InstructorsIDs{})

	router := mux.NewRouter()
	router.HandleFunc("/ws/{userId}", handleConnections).Methods("GET")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	natsConn = service.ConnectNATS()
	service.UpdateCourseInfo(db, natsConn)
	service.UpdateCourseStudents(db, natsConn)

	go handleMessages()

	log.Println("http server started on :8000")
	srv.ListenAndServe()
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["userId"]
	userId, err := strconv.ParseInt(user, 10, 64)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	activeUsers[ws] = userId
	onlineUsers[userId] = true
	go service.SendAllCourseConversation(userId, ws, db)
	go service.SendAllConversations(userId, ws, db)

	profile := service.GetProfile(userId, db)
	if !(profile.Id > 0) {
		profileWrapper, err := service.GetProfileFromNATS(userId, natsConn)
		if err != nil {
			log.Println("Unable to get user profile from userProfile.GetProfile")
			msg := "No data from userProfile.GetProfile"
			msgB, _ := json.Marshal(msg)
			ws.WriteMessage(1, msgB)
			delete(activeUsers, ws)
		} else {
			profile = service.SaveProfileToDB(profileWrapper.Profile, db)
		}
	}
	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(activeUsers, ws)
			break
		}
		if messageType == 1 {
			var message model.Message
			err = json.Unmarshal(data, &message)
			if err != nil {
				log.Printf("error: %v", err)
				delete(activeUsers, ws)
				break
			}
			message.AuthorId = userId
			if !(message.CourseId > 0 || message.RecipientId > 0) {
				var search model.Search
				err = json.Unmarshal(data, &search)
				if err != nil {
					log.Printf("error: %v", err)
					delete(activeUsers, ws)
					break
				}
				if search.Name == "" {
					log.Printf("error: CourseId/RecipientId not present")
					msg := "Missing course_id/recipient_id"
					msgB, _ := json.Marshal(msg)
					ws.WriteMessage(1, msgB)
					delete(activeUsers, ws)
					delete(onlineUsers, message.AuthorId)
					break
				}
				if search.Resource == "course" {
					service.SendSearchResults(userId, search, onlineUsers, activeUsers, db)
				}
			} else if message.Attachment != "" {
				splitAttachment := strings.Split(message.Attachment, ".")
				ext := strings.ToLower(splitAttachment[len(splitAttachment)-1])
				if !model.AllowedExtensions[ext] {
					log.Printf("error: Extension not allowed")
					msg := "Extension not allowed"
					msgB, _ := json.Marshal(msg)
					ws.WriteMessage(1, msgB)
					delete(activeUsers, ws)
					delete(onlineUsers, message.AuthorId)
					break
				}
				go service.SaveMessage(message, natsConn, db)
				continue
			}
			broadcast <- message
		} else {
			fileSize := len(data) / (1024*1024)
			if fileSize > 20 {
				log.Printf("error: File size exceeded")
				msg := "File size must be under 20MB"
				msgB, _ := json.Marshal(msg)
				ws.WriteMessage(1, msgB)
				delete(activeUsers, ws)
				delete(onlineUsers, userId)
				break
			}
			message := service.GetFileFromWS(data, userId, db)
			broadcast <- message
		}
	}
}

func handleMessages() {
	for {
		message := <-broadcast
		if message.RecipientId > 0 {
			service.SendingPrivateMessage(message, onlineUsers, activeUsers)
		} else if message.CourseId > 0 {
			service.SendingGroupMessage(message, activeUsers)
		}
	}
}
