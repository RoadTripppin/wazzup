package migrations

import (
	"github.com/RoadTripppin/wazzup/helpers"
	"github.com/RoadTripppin/wazzup/interfaces"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func createAccounts() {
	db := helpers.ConnectDB()

	users := &[2]interfaces.User{
		{Name: "Martin", Email: "martin@martin.com"},
		{Name: "Michael", Email: "michael@michael.com"},
	}

	for i := 0; i < len(users); i++ {
		// Correct one way
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Name))
		user := &interfaces.User{Name: users[i].Name, Email: users[i].Email, Password: generatedPassword}
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