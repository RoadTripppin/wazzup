package migrations

import (
	"github.com/RoadTripppin/wazzup/helpers"
	"github.com/RoadTripppin/wazzup/interfaces"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func createAccounts() {
	db := helpers.ConnectDB()

	users := &[2]interfaces.User{
		{Username: "Martin", Email: "martin@martin.com"},
		{Username: "Michael", Email: "michael@michael.com"},
	}

	for i := 0; i < len(users); i++ {
		// Correct one way
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))
		user := &interfaces.User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
		db.Create(&user)
	}
	defer db.Close()
}


func Migrate() {
	User := &interfaces.User{}
	db := helpers.ConnectDB()
	db.AutoMigrate(&User)
	defer db.Close()
	
	createAccounts()
}