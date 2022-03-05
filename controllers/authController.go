package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/RoadTripppin/wazzup/helpers"
	"github.com/RoadTripppin/wazzup/models"
	"github.com/RoadTripppin/wazzup/users"
)

func login(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)
	// Handle Login
	var formattedBody models.Login
	err = json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	login := users.Login(formattedBody.Email, formattedBody.Password)
	// Prepare response
	if login["message"] == "all is fine" {
		resp := login
		json.NewEncoder(w).Encode(resp)
		// Handle error in else
	} else {
		resp := models.ErrResponse{Message: "Wrong username or password"}
		json.NewEncoder(w).Encode(resp)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)
	// Handle registration
	var formattedBody models.Register
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
		resp := models.ErrResponse{Message: "Wrong username or password"}
		json.NewEncoder(w).Encode(resp)
	}
}