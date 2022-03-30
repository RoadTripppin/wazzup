package helpers

import (
	"strconv"

	"github.com/RoadTripppin/wazzup/models"
)

func prepareAuthResponse(user *models.User) map[string]interface{} {

	userData := map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"password":   user.Password,
		"profilepic": user.ProfilePic,
	}

	var token = prepareToken(strconv.FormatUint(uint64(user.ID), 10))
	var response = map[string]interface{}{"message": "all is fine"}
	response["token"] = token
	response["user"] = userData

	return response
}

func prepareGetUserResponse(user *models.User) map[string]interface{} {
	userData := map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"profilepic": user.ProfilePic,
	}

	var response = map[string]interface{}{
		"message": "all is fine",
		"user":    userData,
	}

	return response
}
