package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RoadTripppin/wazzup/controllers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func StartApi() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("SERVER_PORT")
	router := mux.NewRouter()
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	fmt.Println("App is working on port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
