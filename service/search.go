package service

import (
	"messaging-service/model"
	"github.com/jinzhu/gorm"
	"messaging-service/database"
)

func SearchProfile(search model.Search, db *gorm.DB) []model.Profile {
	return database.GetProfileByName(search.Name, db)
}

func SearchCourse(search model.Search, db *gorm.DB) []model.DBCourse {
	return database.GetCoursesByField("name", search.Name, db)
}
