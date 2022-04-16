package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

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
		// resp := models.ErrResponse{Message: "Incorrect Details"}
		resp := registerUser
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(resp)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.Register
	_ = json.NewDecoder(r.Body).Decode(&user)

	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		resp := models.ErrResponse{Message: "No Auth token found."}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
	} else {
		splitToken := strings.Split(reqToken, " ")
		reqToken = splitToken[1]

		updatedUser := helpers.UpdateUser(reqToken, &user)
		if updatedUser["message"] == "all is fine" {
			resp := updatedUser

			resp["message"] = "User Updated Successfully"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)

		} else {
			resp := models.ErrResponse{Message: updatedUser["message"].(string)}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(resp)
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body models.Register
	_ = json.NewDecoder(r.Body).Decode(&body)

	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, " ")
	reqToken = splitToken[1]

	deletedUser := helpers.DeleteUser(reqToken, &body)
	if deletedUser["message"] == "all is fine" {
		resp := deletedUser
		resp["message"] = "User Deleted Successfully"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := models.ErrResponse{Message: deletedUser["message"].(string)}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
	}
}
