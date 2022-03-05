package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/RoadTripppin/wazzup/helpers"
	"github.com/RoadTripppin/wazzup/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var login models.Login
	_ = json.NewDecoder(r.Body).Decode(&login)

	loginUser := helpers.Login(login.Email, login.Password)

	if loginUser["message"] == "all is fine" {
		resp := loginUser
		resp["message"] = "User Login Successful"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := models.ErrResponse{Message: "Incorrect Credentials"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var register models.Register
	_ = json.NewDecoder(r.Body).Decode(&register)

	registerUser := helpers.Register(register.Name, register.Email, register.Password, register.ProfilePic)
	// Prepare response
	if registerUser["message"] == "all is fine" {
		resp := registerUser
		resp["message"] = "User registered succesfully"
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
		// Handle error in else
	} else {
		resp := models.ErrResponse{Message: "Incorrect Details"}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(resp)
	}
}
