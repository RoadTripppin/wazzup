package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/RoadTripppin/wazzup/helpers"
	"github.com/RoadTripppin/wazzup/users"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Login struct {
	Email string
	Password string
}

type Register struct {
	Name string
	Email string
	Password string
	ProfilePic string
}

type ErrResponse struct {
	Message string
}

func login(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)
	// Handle Login
	var formattedBody Login
	err = json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	login := users.Login(formattedBody.Email, formattedBody.Password)
	// Prepare response
	if login["message"] == "all is fine" {
		resp := login
		json.NewEncoder(w).Encode(resp)
		// Handle error in else
	} else {
		resp := ErrResponse{Message: "Wrong username or password"}
		json.NewEncoder(w).Encode(resp)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)
	// Handle registration
	var formattedBody Register
	err = json.Unmarshal(body, &formattedBody)
	log.Println(formattedBody)
	helpers.HandleErr(err)
	register := users.Register(formattedBody.Name, formattedBody.Email, formattedBody.Password)
	// Prepare response
	log.Println(register)
	if register["message"] == "all is fine" {
		resp := register
		json.NewEncoder(w).Encode(resp)
		// Handle error in else
	} else {
		resp := ErrResponse{Message: "Wrong username or password"}
		json.NewEncoder(w).Encode(resp)
	}
}


func StartApi() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("SERVER_PORT")
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	fmt.Println("App is working on port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
	
}