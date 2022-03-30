package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/RoadTripppin/wazzup/helpers"
	"github.com/RoadTripppin/wazzup/models"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	email := "asd"
	userResp := helpers.GetUser(email)

	if userResp["message"] == "all is fine" {
		resp := userResp
		resp["message"] = "Found User Successfully"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := models.ErrResponse{Message: "Mo User Found"}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
	}
}
