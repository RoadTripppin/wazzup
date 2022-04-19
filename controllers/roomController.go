package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RoadTripppin/wazzup/helpers"
	"github.com/RoadTripppin/wazzup/models"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body models.GetRoom

	_ = json.NewDecoder(r.Body).Decode(&body)

	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		resp := models.ErrResponse{Message: "No Auth token found."}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
	} else {
		splitToken := strings.Split(reqToken, " ")
		reqToken = splitToken[1]

		messages := helpers.GetMessages(reqToken, body.RoomID)

		if messages["message"] == "all is fine" {
			resp := messages
			resp["message"] = "Room Found"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := models.ErrResponse{Message: "No User Found"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
		}
	}
}
