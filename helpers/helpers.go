package helpers

import (
	"fmt"
	"log"
	"regexp"

	"os"

	"github.com/RoadTripppin/wazzup/models"
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

const projectDirName = "wazzup"

func LoadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	fmt.Println(string(rootPath))
	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatalf("Error loading .env file !")
	}
}

func ConnectDB() *gorm.DB {
	LoadEnv()

	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")
	db_password := os.Getenv("DB_PASSWORD")
	db, err := gorm.Open("postgres", "host="+db_host+" port=5432 user="+db_user+" password="+db_password+" dbname="+db_name+" sslmode=require")
	HandleErr(err)
	return db
}

func Validation(values []models.Validation) (bool, string) {
	name := regexp.MustCompile(`^[a-zA-Z]+(\s[a-zA-Z]+)?$`)
	email := regexp.MustCompile(`^[A-Za-z0-9]+[@]+[A-Za-z0-9]+[.]+[A-Za-z]+$`)

	for i := 0; i < len(values); i++ {
		switch values[i].Valid {
		case "name":
			if !name.MatchString(values[i].Value) {
				return false, "name"
			}
		case "email":
			if !email.MatchString(values[i].Value) {
				return false, "email"
			}
		case "password":
			if len(values[i].Value) < 5 {
				return false, "password"
			}
		}
	}
	return true, ""
}
