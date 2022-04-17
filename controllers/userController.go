package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RoadTripppin/wazzup/helpers"
	"github.com/RoadTripppin/wazzup/models"
)

func SearchUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body models.SearchBody

	_ = json.NewDecoder(r.Body).Decode(&body)

	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		resp := models.ErrResponse{Message: "No Auth token found."}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
	} else {
		splitToken := strings.Split(reqToken, " ")
		reqToken = splitToken[1]

		searchedUser := helpers.SearchUser(reqToken, body.Querystring)

		if searchedUser["message"] == "all is fine" {
			resp := searchedUser
			resp["message"] = "User Found"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := models.ErrResponse{Message: "No User Found"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
		}
	}

	// userResp := helpers.GetUser(email)

	// if userResp["message"] == "all is fine" {
	// 	resp := userResp
	// 	resp["message"] = "Found User Successfully"
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(resp)
	// } else {
	// 	resp := models.ErrResponse{Message: "Mo User Found"}
	// 	w.WriteHeader(http.StatusNotFound)
	// 	json.NewEncoder(w).Encode(resp)
	// }
}
