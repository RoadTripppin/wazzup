package helpers

import (
	"log"
	"regexp"

	"os"

	"github.com/RoadTripppin/wazzup/interfaces"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)


func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	HandleErr(err)

	return string(hashed)
}

func ConnectDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")
	db_password := os.Getenv("DB_PASSWORD")
	// db, err := gorm.Open("postgres", "host=localhost port=5432 user=" + db_user + " dbname=" + db_name + " sslmode=disable")
	db, err := gorm.Open("postgres", "host=" + db_host + " port=5432 user=" + db_user + " password=" + db_password +" dbname=" + db_name + " sslmode=require")
	HandleErr(err)
	return db
}

func Validation(values []interfaces.Validation) bool{
	name := regexp.MustCompile(`^([A-Za-z]{2,})+$`)
	email := regexp.MustCompile(`^[A-Za-z0-9]+[@]+[A-Za-z0-9]+[.]+[A-Za-z]+$`)

	for i := 0; i < len(values); i++ {
		switch values[i].Valid {
			case "name":
				if !name.MatchString(values[i].Value) {
					return false
				}
			case "email":
				if !email.MatchString(values[i].Value) {
					return false
				}
			case "password":
				if len(values[i].Value) < 5 {
					return false
				}
		}
	}
	return true
}