package models

import "github.com/jinzhu/gorm"

type Register struct {
	gorm.Model
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ProfilePic string `json:"profilepic"`
}

type Login struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SearchBody struct {
	Querystring string
}

type GetRoom struct {
	RoomID string `json:"roomid"`
}

type RoomData struct {
	ID       string
	Name     string
	Messages string
}

type User struct {
	gorm.Model
	Name       string
	Email      string
	Password   string
	ProfilePic string
}

type Validation struct {
	Value string
	Valid string
}
